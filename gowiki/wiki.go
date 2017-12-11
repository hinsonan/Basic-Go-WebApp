package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

//used to define page data object
type Page struct {
	Title string
	Body  []byte
}

//allows us to save a page with a title and body
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//loads the name of the file
func loadPage(title string) (*Page, error) {
	//adds the .txt extension to the file
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	//returns the  error if it exist
	if err != nil {
		return nil, err
	}
	//returns the new page
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	//if there is not a page then we want to send them to the edit so they can add content
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

/*
The page title (provided in the URL) and the form's only field, Body, are stored in a new Page.
 The save() method is then called to write the data to a file, and the client is redirected to the /view/ page.

The value returned by FormValue is of type string.
We must convert that value to []byte before it will fit into the Page struct.
We use []byte(body) to perform the conversion.
*/
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//writes to the specific html template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func main() {
	/*
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.save()
		p2, _ := loadPage("TestPage")
		fmt.Println(string(p2.Body))
	*/
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
