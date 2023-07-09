package main

import (
	"errors"
	"fmt"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	"github.com/b-charles/pigs/smartconfig"
)

type meteoMessage struct {
	today string

	sunrise  string
	sunset   string
	sundelta int

	morningT   int
	morningW   int
	afternoonT int
	afternoonW int
	eveningT   int
	eveningW   int
	nightT     int
	nightW     int
}

func (self *meteoMessage) format() *Message {

	var (
		sunrise  string
		sunset   string
		sundelta string
	)
	if self.sunrise == "" || self.sunset == "" {
		sunrise = "Tôt"
		sunset = "Tard"
		sundelta = "Comme hier"
	} else {
		sunrise = self.sunrise
		sunset = self.sunset
		sundelta = fmt.Sprintf("%+dmin", self.sundelta)
	}

	morningT, morningW := getForecast(self.morningT, self.morningW, 2, "Frais mais pas trop")
	afternoonT, afternoonW := getForecast(self.afternoonT, self.afternoonW, 99, "Chaud bouillant")
	eveningT, eveningW := getForecast(self.eveningT, self.eveningW, 21, "Tièdasse")
	nightT, nightW := getForecast(self.nightT, self.nightW, -9999999, "Sibérique")

	msg := fmt.Sprintf("Météo Paris &#2947; %s &#2947; "+
		"Lever du soleil: %s Coucher du soleil: %s (%s) &#2947; "+
		"Temps de la journée &#2947; "+
		"Matin: %d° %s &#2947; "+
		"Après-midi : %d° %s &#2947; "+
		"Soirée : %d° %s &#2947; "+
		"Nuit : %d° %s &#2947; "+
		"Bonne journée !",
		self.today,
		sunrise, sunset, sundelta,
		morningT, morningW,
		afternoonT, afternoonW,
		eveningT, eveningW,
		nightT, nightW)

	return &Message{Title: "Météo Paris", Body: msg}

}

type MeteoConfig struct {
	AccessToken string
}

type Meteo struct {
	Today      *Today       `inject:""`
	Config     *MeteoConfig `inject:""`
	HttpClient HttpClient   `inject:""`
}

func (self *Meteo) Order() int {
	return 10
}

func (self *Meteo) Name() string {
	return "meteopnm"
}

const (
	METEO_API_BASE string = "https://api.meteo-concept.com/api"
	INSEE          string = "75056"
)

func (self *Meteo) Message() (*Message, error) {

	msg := &meteoMessage{}

	msg.today = self.Today.Get()

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", self.Config.AccessToken),
	}

	ephemeride := self.HttpClient.Get(
		fmt.Sprintf("%s%s?insee=%s", METEO_API_BASE, "/ephemeride/0", INSEE),
		headers)
	forecast := self.HttpClient.Get(
		fmt.Sprintf("%s%s?insee=%s", METEO_API_BASE, "/forecast/daily/periods", INSEE),
		headers)

	var ephemerideErr error = nil
	if resp := <-ephemeride; resp.Err != nil {
		ephemerideErr = resp.Err
	} else {
		if node, err := json.ParseString(resp.Body); err != nil {
			ephemerideErr = err
		} else {

			eph := node.GetMember("ephemeride")
			msg.sunrise = eph.GetMember("sunrise").AsString()
			msg.sunset = eph.GetMember("sunset").AsString()
			msg.sundelta = eph.GetMember("diff_duration_day").AsInt()

		}
	}

	var forecastErr error = nil
	if resp := <-forecast; resp.Err != nil {
		forecastErr = resp.Err
	} else {
		if node, err := json.ParseString(resp.Body); err != nil {
			forecastErr = err
		} else {

			today := node.GetMember("forecast").GetElement(0)
			msg.morningT = today.GetElement(1).GetMember("temp2m").AsInt()
			msg.morningW = today.GetElement(1).GetMember("weather").AsInt()
			msg.afternoonT = today.GetElement(2).GetMember("temp2m").AsInt()
			msg.afternoonW = today.GetElement(2).GetMember("weather").AsInt()
			msg.eveningT = today.GetElement(3).GetMember("temp2m").AsInt()
			msg.eveningW = today.GetElement(3).GetMember("weather").AsInt()

			tomorrow := node.GetMember("forecast").GetElement(1)
			msg.nightT = tomorrow.GetElement(0).GetMember("temp2m").AsInt()
			msg.nightW = tomorrow.GetElement(0).GetMember("weather").AsInt()

		}
	}

	return msg.format(), errors.Join(ephemerideErr, forecastErr)

}

func init() {
	smartconfig.Configure("meteo", &MeteoConfig{})
	ioc.Put(&Meteo{}, func(Service) {})
}

func getForecast(temp int, weather int, defTemp int, defWeather string) (int, string) {

	switch weather {
	case 0:
		return temp, "Soleil"
	case 1:
		return temp, "Peu nuageux"
	case 2:
		return temp, "Ciel voilé"
	case 3:
		return temp, "Nuageux"
	case 4:
		return temp, "Très nuageux"
	case 5:
		return temp, "Couvert"
	case 6:
		return temp, "Brouillard"
	case 7:
		return temp, "Brouillard givrant"
	case 10:
		return temp, "Pluie faible"
	case 11:
		return temp, "Pluie modérée"
	case 12:
		return temp, "Pluie forte"
	case 13:
		return temp, "Pluie faible verglaçante"
	case 14:
		return temp, "Pluie modérée verglaçante"
	case 15:
		return temp, "Pluie forte verglaçante"
	case 16:
		return temp, "Bruine"
	case 20:
		return temp, "Neige faible"
	case 21:
		return temp, "Neige modérée"
	case 22:
		return temp, "Neige forte"
	case 30:
		return temp, "Pluie et neige mêlées faibles"
	case 31:
		return temp, "Pluie et neige mêlées modérées"
	case 32:
		return temp, "Pluie et neige mêlées fortes"
	case 40:
		return temp, "Averses de pluie locales et faibles"
	case 41:
		return temp, "Averses de pluie locales"
	case 42:
		return temp, "Averses locales et fortes"
	case 43:
		return temp, "Averses de pluie faibles"
	case 44:
		return temp, "Averses de pluie"
	case 45:
		return temp, "Averses de pluie fortes"
	case 46:
		return temp, "Averses de pluie faibles et fréquentes"
	case 47:
		return temp, "Averses de pluie fréquentes"
	case 48:
		return temp, "Averses de pluie fortes et fréquentes"
	case 60:
		return temp, "Averses de neige localisées et faibles"
	case 61:
		return temp, "Averses de neige localisées"
	case 62:
		return temp, "Averses de neige localisées et fortes"
	case 63:
		return temp, "Averses de neige faibles"
	case 64:
		return temp, "Averses de neige"
	case 65:
		return temp, "Averses de neige fortes"
	case 66:
		return temp, "Averses de neige faibles et fréquentes"
	case 67:
		return temp, "Averses de neige fréquentes"
	case 68:
		return temp, "Averses de neige fortes et fréquentes"
	case 70:
		return temp, "Averses de pluie et neige mêlées localisées et faibles"
	case 71:
		return temp, "Averses de pluie et neige mêlées localisées"
	case 72:
		return temp, "Averses de pluie et neige mêlées localisées et fortes"
	case 73:
		return temp, "Averses de pluie et neige mêlées faibles"
	case 74:
		return temp, "Averses de pluie et neige mêlées"
	case 75:
		return temp, "Averses de pluie et neige mêlées fortes"
	case 76:
		return temp, "Averses de pluie et neige mêlées faibles et nombreuses"
	case 77:
		return temp, "Averses de pluie et neige mêlées fréquentes"
	case 78:
		return temp, "Averses de pluie et neige mêlées fortes et fréquentes"
	case 100:
		return temp, "Orages faibles et locaux"
	case 101:
		return temp, "Orages locaux"
	case 102:
		return temp, "Orages fort et locaux"
	case 103:
		return temp, "Orages faibles"
	case 104:
		return temp, "Orages"
	case 105:
		return temp, "Orages forts"
	case 106:
		return temp, "Orages faibles et fréquents"
	case 107:
		return temp, "Orages fréquents"
	case 108:
		return temp, "Orages forts et fréquents"
	case 120:
		return temp, "Orages faibles et locaux de neige ou grésil"
	case 121:
		return temp, "Orages locaux de neige ou grésil"
	case 122:
		return temp, "Orages locaux de neige ou grésil"
	case 123:
		return temp, "Orages faibles de neige ou grésil"
	case 124:
		return temp, "Orages de neige ou grésil"
	case 125:
		return temp, "Orages de neige ou grésil"
	case 126:
		return temp, "Orages faibles et fréquents de neige ou grésil"
	case 127:
		return temp, "Orages fréquents de neige ou grésil"
	case 128:
		return temp, "Orages fréquents de neige ou grésil"
	case 130:
		return temp, "Orages faibles et locaux de pluie et neige mêlées ou grésil"
	case 131:
		return temp, "Orages locaux de pluie et neige mêlées ou grésil"
	case 132:
		return temp, "Orages fort et locaux de pluie et neige mêlées ou grésil"
	case 133:
		return temp, "Orages faibles de pluie et neige mêlées ou grésil"
	case 134:
		return temp, "Orages de pluie et neige mêlées ou grésil"
	case 135:
		return temp, "Orages forts de pluie et neige mêlées ou grésil"
	case 136:
		return temp, "Orages faibles et fréquents de pluie et neige mêlées ou grésil"
	case 137:
		return temp, "Orages fréquents de pluie et neige mêlées ou grésil"
	case 138:
		return temp, "Orages forts et fréquents de pluie et neige mêlées ou grésil"
	case 140:
		return temp, "Pluies orageuses"
	case 141:
		return temp, "Pluie et neige mêlées à caractère orageux"
	case 142:
		return temp, "Neige à caractère orageux"
	case 210:
		return temp, "Pluie faible intermittente"
	case 211:
		return temp, "Pluie modérée intermittente"
	case 212:
		return temp, "Pluie forte intermittente"
	case 220:
		return temp, "Neige faible intermittente"
	case 221:
		return temp, "Neige modérée intermittente"
	case 222:
		return temp, "Neige forte intermittente"
	case 230:
		return temp, "Pluie et neige mêlées"
	case 231:
		return temp, "Pluie et neige mêlées"
	case 232:
		return temp, "Pluie et neige mêlées"
	case 235:
		return temp, "Averses de grêle"
	default:
		return defTemp, defWeather
	}

}
