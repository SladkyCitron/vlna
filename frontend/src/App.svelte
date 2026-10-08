<script lang="ts">
  import TitleBar from "$lib/components/TitleBar.svelte";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import Player from "$lib/components/Player.svelte";
  import StationDetails from "$lib/components/StationDetails.svelte";
  import Explore from "$lib/components/views/Explore.svelte";
  import Favorites from "$lib/components/views/Favorites.svelte";
  import Search from "$lib/components/views/Search.svelte";
  import Settings from "$lib/components/views/Settings.svelte";
  import { activeView } from "$lib/stores/nav";
  import { theme } from "$lib/stores/theme";

  $: if (typeof document !== "undefined") {
    document.documentElement.setAttribute("data-theme", $theme ?? "light");
  }
</script>

<div
  class="bg-background text-foreground grid h-screen w-screen grid-cols-[240px_1fr] grid-rows-[auto_1fr_auto] overflow-hidden select-none"
>
  <TitleBar />

  <div class="col-span-2 flex min-h-0 overflow-hidden">
    <Sidebar />

    <main class="flex-1 overflow-y-auto p-4">
      {#if $activeView === "explore"}
        <Explore />
      {:else if $activeView === "search"}
        <Search />
      {:else if $activeView === "favorites"}
        <Favorites />
      {:else if $activeView === "settings"}
        <Settings />
      {/if}
    </main>
  </div>

  <Player />
</div>

<StationDetails />
