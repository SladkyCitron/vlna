<script lang="ts">
  import {
    Heart,
    Pause,
    Play,
    Radio,
    Volume2,
    LoaderCircle,
    TriangleAlert,
  } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button";
  import * as m from "$lib/paraglide/messages.js";
  import {
    currentStation,
    status,
    currentSongMetadata,
    volume,
    togglePlay,
    setPlayerVolume,
  } from "$lib/stores/player";
  import { toggleFavorite, favorites } from "$lib/stores/favorites";
  import { getLocalizedCountryName } from "$lib/utils";

  function handleVolumeChange(e: Event) {
    const target = e.target as HTMLInputElement;
    setPlayerVolume(Number(target.value));
  }

  function formatSongMetadata() {
    // format the song metadata to display "Artist - Title"
    // and make Hatsune Miku's name bold (easter egg)
    if ($currentSongMetadata?.artist) {
      return `${$currentSongMetadata.artist} - ${$currentSongMetadata.title}`.replace(
        /hatsune.miku/i,
        "<strong>Hatsune Miku</strong>"
      );
    } else if ($currentSongMetadata?.title) {
      return $currentSongMetadata.title.replace(
        /hatsune.miku/i,
        "<strong>Hatsune Miku</strong>"
      );
    } else {
      return "";
    }
  }
</script>

<footer
  class="border-border bg-card col-span-2 row-start-3 flex min-h-20 w-full items-center gap-4 border-t px-4 py-3 select-none"
>
  <div class="flex min-w-0 flex-1 items-center gap-3">
    <div
      class="bg-muted text-muted-foreground flex size-12 shrink-0 items-center justify-center rounded-md"
    >
      {#if $currentStation?.favicon}
        <img
          src={$currentStation.favicon}
          alt={$currentStation.name}
          class="size-full object-cover"
        />
      {:else}
        <Radio class="size-6" />
      {/if}
    </div>
    <div class="min-w-0">
      <p class="text-foreground truncate text-sm font-semibold">
        {$currentStation ? $currentStation.name : m.noStationSelected()}
      </p>
      <p class="text-muted-foreground truncate text-xs">
        {#if $currentSongMetadata?.title}
          <span class="text-primary font-medium">
            {@html formatSongMetadata()}
          </span>
        {:else if $currentStation}
          {getLocalizedCountryName($currentStation.countrycode)}
        {:else}
          {m.chooseStation()}
        {/if}
      </p>
    </div>
    {#if $currentStation}
      <Button
        variant="ghost"
        size="icon"
        class="text-muted-foreground ml-1 shrink-0"
        onclick={() => toggleFavorite($currentStation!)}
      >
        <Heart
          class="size-4"
          fill={($favorites ?? []).some(
            (f) => f.stationuuid === $currentStation?.stationuuid
          )
            ? "currentColor"
            : "none"}
        />
      </Button>
    {/if}
  </div>

  <div class="flex shrink-0 items-center">
    <Button
      variant="secondary"
      size="icon-lg"
      class="rounded-full"
      disabled={!$currentStation}
      onclick={togglePlay}
    >
      {#if $status === "Playing"}
        <Pause class="size-5 fill-current" />
      {:else if $status === "Paused" || $status === "Stopped"}
        <Play class="size-5 fill-current" />
      {:else if $status === "Loading"}
        <LoaderCircle class="size-5 animate-spin" />
      {:else}
        <TriangleAlert class="size-5 fill-current" />
      {/if}
    </Button>
  </div>

  <div class="flex min-w-0 flex-1 items-center justify-end gap-2">
    <Volume2 class="text-muted-foreground size-4 shrink-0" />
    <input
      class="accent-primary hidden h-1 w-24 cursor-pointer sm:block"
      type="range"
      min="0"
      max="100"
      value={$volume}
      oninput={handleVolumeChange}
    />
    <span
      class="text-muted-foreground hidden w-8 text-right text-xs tabular-nums md:block"
    >
      {$volume}%
    </span>
  </div>
</footer>
