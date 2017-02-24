package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/assets"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/go-openapi/runtime"
	http_transport "github.com/go-openapi/runtime/client"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync/atomic"
)

type Checker struct {
	// Имя текущей проверки.
	Name string
	// Описание текущей проверки.
	Description string
	// Функция для текущей проверки.
	FnCheck func(c *client.Forum)
	// Тесты, без которых проверка не имеет смысл.
	Deps []string
}

var s_templateUid int32 = 0

type CheckerByName []Checker

func (a CheckerByName) Len() int           { return len(a) }
func (a CheckerByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CheckerByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type CheckerTransport struct {
	t      runtime.ClientTransport
	report *Report
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := NewValidator(operation.Context, self.report)
	operation.Client = &http.Client{Transport: tracker}
	return self.t.Submit(operation)
}

func Checkpoint(c *client.Forum, message string) bool {
	return c.Transport.(*CheckerTransport).report.Checkpoint(message)
}

var registeredChecks []Checker

func Register(checker Checker) {
	registeredChecks = append(registeredChecks, checker)
}

func RunCheck(check Checker, report *Report, cfg *client.TransportConfig) {
	report.Result = Success
	transport := http_transport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	defer func() {
		if r := recover(); r != nil {
			report.AddError(r)
		}
	}()
	check.FnCheck(client.New(&CheckerTransport{transport, report}, nil))
}

func SortedChecks() []Checker {
	pending := map[string]Checker{}
	for _, check := range registeredChecks {
		if _, ok := pending[check.Name]; ok {
			log.Fatal("Found duplicate check:", check.Name)
		}
		pending[check.Name] = check
	}

	result := []Checker{}
	added := map[string]bool{}
	for len(pending) > 0 {
		batch := []Checker{}
		// Found ready tasks
		for _, item := range pending {
			ready := true
			for _, dep := range item.Deps {
				if !added[dep] {
					ready = false
					break
				}
			}
			if ready {
				batch = append(batch, item)
			}
		}
		if len(batch) == 0 {
			log.Fatal("Can't found dependencies for tasks:", pending)
		}
		// Sort batch by name
		sort.Sort(CheckerByName(batch))
		// Add ready tasks to result
		for _, item := range batch {
			added[item.Name] = true
			delete(pending, item.Name)
		}
		result = append(result, batch...)
	}

	return result
}

func templateUid() string {
	return fmt.Sprintf("i%d", atomic.AddInt32(&s_templateUid, 1))
}

func templateAsset(outer, name string) template.HTML {
	data, err := assets.Asset(name)
	tag := strings.SplitN(outer, " ", 2)[0]
	if err != nil {
		panic(err)
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", outer, string(data), tag))
}

func reportTemplate() *template.Template {
	data, err := assets.Asset("template.html")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.
		New("template.html").
		Funcs(template.FuncMap{
			"uid":   templateUid,
			"asset": templateAsset,
		}).
		Parse(string(data))
	if err != nil {
		panic(err)
	}
	return tmpl
}

func Run(url *url.URL, keep bool) int {
	total := 0
	failed := 0
	skipped := 0
	broken := map[string]bool{}

	tpl := reportTemplate()
	cfg := client.DefaultTransportConfig().WithHost(url.Host).WithSchemes([]string{url.Scheme}).WithBasePath(url.Path)
	reports := []*Report{}
	for _, check := range SortedChecks() {
		log.Printf("=== RUN:  %s", check.Name)
		report := Report{
			Checker: check,
		}
		for _, dep := range check.Deps {
			if broken[dep] {
				report.Skip(dep)
			}
		}
		if report.Result != Skipped {
			RunCheck(check, &report, cfg)
		}
		if report.Result != Success {
			report.Show()
		}
		var result string
		total++
		switch report.Result {
		case Skipped:
			broken[check.Name] = true
			skipped++
			result = "SKIPPED"
		case Success:
			result = "OK"
		default:
			broken[check.Name] = true
			failed++
			result = "FAILED"
		}
		log.Printf("--- DONE: %s (%s)", check.Name, result)
		reports = append(reports, &report)
		if failed > 0 && !keep {
			break
		}
	}

	file, err := os.Create("report.html")
	defer file.Close()
	err = tpl.Execute(file, struct {
		Reports []*Report
	}{
		Reports: reports,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("RESULT: %d total, %d skipped, %d failed)", total, skipped, failed)
	return failed
}
