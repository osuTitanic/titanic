package wiki

import (
	"fmt"
	"time"
)

var localizedDays = map[string][]string{
	"en": {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	"de": {"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
	"ru": {"Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"},
	"pl": {"Niedziela", "Poniedziałek", "Wtorek", "Środa", "Czwartek", "Piątek", "Sobota"},
	"fi": {"Sunnuntai", "Maanantai", "Tiistai", "Keskiviikko", "Torstai", "Perjantai", "Lauantai"},
	"et": {"Pühapäev", "Esmaspäev", "Teisipäev", "Kolmapäev", "Neljapäev", "Reede", "Laupäev"},
	"fr": {"Dimanche", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"},
	"es": {"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"},
}

var localizedMonths = map[string][]string{
	"en": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	"de": {"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
	"ru": {"Января", "Февраля", "Марта", "Апреля", "Мая", "Июня", "Июля", "Августа", "Сентября", "Октября", "Ноября", "Декабря"},
	"pl": {"Stycznia", "Lutego", "Marca", "Kwietnia", "Maja", "Czerwca", "Lipca", "Sierpnia", "Września", "Października", "Listopada", "Grudnia"},
	"fi": {"Tammikuuta", "Helmikuuta", "Maaliskuuta", "Huhtikuuta", "Toukokuuta", "Kesäkuuta", "Heinäkuuta", "Elokuuta", "Syyskuuta", "Lokakuuta", "Marraskuuta", "Joulukuuta"},
	"et": {"Jaanuar", "Veebruar", "Märts", "Aprill", "Mai", "Juuni", "Juuli", "August", "September", "Oktoober", "November", "Detsember"},
	"fr": {"Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre"},
	"es": {"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"},
}

// CurrentDate returns the current date in a localized format based on the provided language code.
// e.g. "en" -> "Monday, January 1, 2024", "de" -> "Montag, 1. Januar 2024", ...
func CurrentDate(language string, now time.Time) string {
	language = NormalizeLanguage(language)

	days, ok := localizedDays[language]
	if !ok {
		days = localizedDays[DefaultLanguage]
	}
	months, ok := localizedMonths[language]
	if !ok {
		months = localizedMonths[DefaultLanguage]
	}

	day := days[int(now.Weekday())]
	month := months[int(now.Month())-1]

	switch language {
	case "en":
		return fmt.Sprintf("%s, %s %d, %d", day, month, now.Day(), now.Year())
	case "de":
		return fmt.Sprintf("%s, %d. %s %d", day, now.Day(), month, now.Year())
	default:
		return fmt.Sprintf("%s, %d %s %d", day, now.Day(), month, now.Year())
	}
}
