<script lang="ts">
  import * as Item from "$lib/components/ui/item";
  import * as ButtonGroup from "$lib/components/ui/button-group";
  import * as Field from "$lib/components/ui/field";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { LoaderCircle, Radio, Heart, Play } from "@lucide/svelte";
  import { onMount } from "svelte";
  import {
    type Stations,
    type Station,
    StationService,
  } from "$bindings/github.com/SladkyCitron/vlna/service";
  import { openStationDetails } from "$lib/stores/stationDetails";
  import {
    favorites,
    loadFavorites,
    toggleFavorite,
  } from "$lib/stores/favorites";
  import { playStation } from "$lib/stores/player";
  import { getLocalizedCountryName } from "$lib/utils";
  import * as m from "$lib/paraglide/messages.js";

  let stations = $state<Stations | null>(null);
  let loading = $state(false);
  let query = $state("");

  function stopStationDetails(event: MouseEvent) {
    event.stopPropagation();
  }

  function toggleStationFavorite(event: MouseEvent, station: Station) {
    stopStationDetails(event);
    toggleFavorite(station);
  }

  function play(event: MouseEvent, station: Station) {
    stopStationDetails(event);
    playStation(station);
  }

  function searchStations() {
    if (query === "") {
      stations = null;
      return;
    }
    loading = true;
    stations = null;
    StationService.GetStationsByName(query)
      .then((result) => {
        stations = result;
      })
      .catch((error) => {
        console.error("Failed to search stations:", error);
      })
      .finally(() => {
        loading = false;
      });
  }

  onMount(() => {
    void loadFavorites();
  });
</script>

<div>
  <h1 class="pb-4 text-xl font-bold">{m.search()}</h1>
  <Field.Field orientation="horizontal" class="mb-4">
    <Input type="search" placeholder={m.search()} bind:value={query} />
    <Button onclick={searchStations}>{m.search()}...</Button>
  </Field.Field>
  {#if stations}
    <div class="flex flex-col gap-4">
      {#each stations as station}
        <Item.Root variant="outline">
          {#snippet child({ props })}
            <a href="#/" onclick={() => openStationDetails(station)} {...props}>
              <Item.Media variant="image">
                {#if station.favicon}
                  <img
                    src={station.favicon}
                    alt={station.name}
                    class="size-16"
                  />
                {:else}
                  <Radio class="size-8" aria-label={station.name} />
                {/if}
              </Item.Media>
              <Item.Content>
                <Item.Title>{station.name}</Item.Title>
                <Item.Description>
                  {getLocalizedCountryName(station.countrycode)}
                  {#if station.tags}
                    | {station.tags
                      .split(",")
                      .slice(0, 4) // limit to first 5 tags only
                      .join(", ")}
                  {/if}
                </Item.Description>
              </Item.Content>
              <Item.Actions>
                <ButtonGroup.Root>
                  <Button
                    variant="outline"
                    onclick={(event) => toggleStationFavorite(event, station)}
                  >
                    <Heart
                      fill={($favorites ?? []).some(
                        (favorite) =>
                          favorite.stationuuid === station.stationuuid
                      )
                        ? "currentColor"
                        : "none"}
                    />
                  </Button>
                  <Button
                    variant="outline"
                    onclick={(event) => play(event, station)}
                  >
                    <Play />
                  </Button>
                </ButtonGroup.Root>
              </Item.Actions>
            </a>
          {/snippet}
        </Item.Root>
      {/each}
    </div>
  {:else}
    {#if loading}
      <div class="flex items-center justify-center">
        <LoaderCircle class="h-8 w-8 animate-spin" />
        <p class="ml-2">{m.loading()}</p>
      </div>
    {/if}
  {/if}
</div>
