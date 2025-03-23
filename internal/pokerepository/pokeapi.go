package pokerepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zcamcal/pokedexcli/internal/pokecache"
)

type pokeIterator[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type pokeLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokeEncounter struct {
	Encounters []poke `json:"pokemon_encounters"`
}

type poke struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type PokeApi struct {
	cache *pokecache.Cache
}

type pokemoon struct {
	BaseExperience int `json:"base_experience"`
}

func NewPokeApi(interval time.Duration) PokeApi {
	cache := pokecache.NewCache(interval)
	return PokeApi{cache: &cache}
}

func (api *PokeApi) Pokemon(name string) (int, error) {
	path := "https://pokeapi.co/api/v2/pokemon/" + name

	var raw []byte

	var ok bool
	raw, ok = api.cache.Get(path)

	if !ok {
		res, err := http.Get(path)
		if err != nil {
			return 0, err
		}

		defer res.Body.Close()

		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return 0, err
		}

		api.cache.Add(path, raw)
	}

	var pokeResponse pokemoon
	if err := json.Unmarshal(raw, &pokeResponse); err != nil {
		return 0, err
	}

	return pokeResponse.BaseExperience, nil
}

func (api *PokeApi) Encounters(area string) ([]string, error) {
	path := "https://pokeapi.co/api/v2/location-area/" + area

	var pokeResponse pokeEncounter

	var raw []byte

	var ok bool
	raw, ok = api.cache.Get(path)

	if !ok {
		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		api.cache.Add(path, raw)
	}

	if err := json.Unmarshal(raw, &pokeResponse); err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, val := range pokeResponse.Encounters {
		names = append(names, val.Pokemon.Name)
	}

	return names, nil
}

func (api PokeApi) Locations(limit, page int) ([]string, error) {
	if limit <= 0 {
		return nil, errors.New("limit cant be below 0")
	}

	if page <= 0 {
		return nil, errors.New("page cant be below 0")
	}

	var raw []byte
	var ok bool

	path := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=%d&offset=%d", limit, page)
	raw, ok = api.cache.Get(path)
	if !ok {
		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		api.cache.Add(path, raw)
	}

	var pokeResponse pokeIterator[pokeLocation]
	if err := json.Unmarshal(raw, &pokeResponse); err != nil {
		return nil, err
	}

	locations := make([]string, 0)
	for _, val := range pokeResponse.Results {
		locations = append(locations, val.Name)
	}

	return locations, nil
}
