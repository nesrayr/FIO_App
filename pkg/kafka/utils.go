package kafka

import (
	"FIO_App/pkg/dtos"
	"encoding/json"
	"fmt"
	"net/http"
)

func EnrichData(fio FIO) (dtos.PersonDTO, error) {
	name := fio.Name
	age, gender, nationality, err := fetchExternalData(name)
	if err != nil {
		return dtos.PersonDTO{}, err
	}
	return dtos.SetPersonDTO(name, fio.Surname, fio.Patronymic, age, gender, nationality), nil
}

func fetchExternalData(name string) (int, string, string, error) {
	age, err := fetchAge(name)
	if err != nil {
		return 0, "", "", err
	}
	gender, err := fetchGender(name)
	if err != nil {
		return 0, "", "", err
	}
	nationality, err := fetchNationality(name)
	if err != nil {
		return 0, "", "", err
	}

	return age, gender, nationality, nil
}

func fetchAge(name string) (int, error) {
	ageApiUrl := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(ageApiUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var ageResponse struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&ageResponse); err != nil {
		return 0, err
	}

	return ageResponse.Age, nil
}

func fetchGender(name string) (string, error) {
	genderApiUrl := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(genderApiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genderResponse struct {
		Count       int     `json:"count"`
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&genderResponse); err != nil {
		return "", err
	}

	return genderResponse.Gender, nil
}

func fetchNationality(name string) (string, error) {
	nationalityApiUrl := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := http.Get(nationalityApiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var nationalityResponse struct {
		Count     int       `json:"count"`
		Name      string    `json:"name"`
		Countries []Country `json:"country"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&nationalityResponse); err != nil {
		return "", err
	}

	return nationalityResponse.Countries[0].Country, nil
}

type Country struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
