package pokerepository

import (
	"encoding/json"
	"github.com/zcamcal/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

type PokeIterator[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type PokeLocation struct {
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

func NewPokeApi(interval time.Duration) PokeApi {
	cache := pokecache.NewCache(interval)
	return PokeApi{cache: &cache}
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
