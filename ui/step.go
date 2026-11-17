package ui

import (
	clientdirectory "github.com/veltylabs/client_directory"
	"github.com/veltylabs/client_directory/kinds"
	"webtyp.com/crud"
	. "webtyp.com/dom"
	h "webtyp.com/html"
	"webtyp.com/input"
	"webtyp.com/lang"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

const (
	textDocumentNumber lang.Text = "Document number"
	textSearchPatient  lang.Text = "Search patient"
	textFullName       lang.Text = "Full name"
	textPhone          lang.Text = "Phone"
	textEmail          lang.Text = "Email"
	textRegisterNew    lang.Text = "Register patient"
	textUseThisPatient lang.Text = "Continue"
)

// PatientStep es el paso «paciente» de una pantalla de otro módulo (p. ej. la reserva de mesón): el personal escribe
// el documento, ve la ficha encontrada o registra al paciente, y la elige. Quien lo usa lo recibe como una interfaz
// mínima suya (`Render(choose)` y `Reset()`); este módulo no importa a nadie.
//
// Ejemplo: step := ui.NewPatientStep(d); el = step.Render(func(id, name string) { /* avanzar */ })
type PatientStep struct {
	caller router.Caller
	ids    model.IDGenerator
	deps   view.BrowserDeps

	document, name, phone, email *SignalString
	errMsg                       *SignalString
	result                       *SignalNodes
	choose                       func(id, name string)
}

// NewPatientStep arma el paso con lo que la app entrega a las pantallas (conexión, ids, países).
func NewPatientStep(d view.BrowserDeps) *PatientStep {
	return &PatientStep{
		caller: d.Caller, ids: d.IDs, deps: d,
		document: NewString(""), name: NewString(""), phone: NewString(""), email: NewString(""),
		errMsg: NewString(""), result: NewNodes(),
	}
}

// Render devuelve el paso; choose se llama con el id y el nombre del paciente elegido.
func (s *PatientStep) Render(choose func(id, name string)) *Element {
	s.choose = choose
	return h.Div().Child(
		h.P().Attr("role", "alert").BindText(s.errMsg),
		field(textDocumentNumber, "patient-document", "text", s.document),
		button(textSearchPatient, s.search),
		h.Div().BindChildren(s.result),
	)
}

// Reset deja el paso limpio para una reserva nueva.
func (s *PatientStep) Reset() {
	s.errMsg.Set("")
	s.result.Set(nil)
	s.document.Set("")
	s.name.Set("")
	s.phone.Set("")
	s.email.Set("")
}

func (s *PatientStep) operation(action model.Action) router.Operation {
	return router.Operation{Module: clientdirectory.ModelName, Entity: (&clientdirectory.Client{}).ModelName(), Action: action}
}

// search busca por el documento escrito; sin coincidencia exacta ofrece registrar.
func (s *PatientStep) search() {
	s.errMsg.Set("")
	s.result.Set(nil)
	doc := s.document.Get()
	want := doc
	if normalized, err := input.NormalizeRUT(doc); err == nil {
		want = normalized
	}
	var found clientdirectory.ClientList
	s.caller.Call(s.operation(model.Read), &crud.ReadArgs{Search: doc, Limit: crud.MaxPage}, &found, func(err error) {
		if err != nil {
			s.errMsg.Set(err.Error())
			return
		}
		for _, p := range found {
			if p.Rut == want || p.Rut == doc {
				s.showCard(*p)
				return
			}
		}
		s.result.Set([]*Element{s.registerForm()})
	})
}

func (s *PatientStep) registerForm() *Element {
	return h.Div().Child(
		field(textFullName, "patient-name", "text", s.name),
		field(textPhone, "patient-phone", "tel", s.phone),
		field(textEmail, "patient-email", "email", s.email),
		button(textRegisterNew, s.register),
	)
}

func (s *PatientStep) register() {
	s.errMsg.Set("")
	p := &clientdirectory.Client{
		Id: s.ids.NewID(), Rut: s.document.Get(), DocumentType: kinds.InitialDocument(s.deps.Countries),
		Name: s.name.Get(), Phone: s.phone.Get(), Email: s.email.Get(), IsActive: true,
	}
	var saved clientdirectory.Client
	s.caller.Call(s.operation(model.Create), &crud.CreateArgs{Records: []model.Model{p}}, &saved, func(err error) {
		if err != nil {
			s.errMsg.Set(err.Error())
			return
		}
		s.showCard(saved)
	})
}

func (s *PatientStep) showCard(p clientdirectory.Client) {
	s.result.Set([]*Element{h.Div().Child(
		h.Strong().Text(p.Name),
		h.P().Text(p.Phone),
		button(textUseThisPatient, func() { s.choose(p.Id, p.Name) }),
	)})
}

func field(label lang.Text, name, kind string, value *SignalString) *Element {
	return h.Label().Text(lang.Translate(label).String()).Child(h.Input(kind).Attr("name", name).Bind(value))
}

func button(label lang.Text, run func()) *Element {
	return h.Button().Attr("type", "button").Text(lang.Translate(label).String()).OnClick(func(Event) { run() })
}
