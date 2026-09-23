<script lang="ts">
  import TitleBar from "$lib/components/TitleBar.svelte";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { activeView } from "$lib/stores/nav";
</script>

<div
  class="dark h-screen w-screen bg-background text-foreground overflow-hidden select-none grid grid-rows-[auto_1fr_auto] grid-cols-[240px_1fr]"
>
  <TitleBar />

  <div class="col-span-2 min-h-0 flex overflow-hidden">
    <Sidebar />

    <main class="flex-1 p-4 overflow-y-auto">
      {#if $activeView === "explore"}
        <div>
          <h1 class="text-xl font-bold">Explore TODO</h1>
        </div>
      {:else if $activeView === "favorites"}
        <div>
          <h1 class="text-xl font-bold">Favorites TODO</h1>
        </div>
      {:else if $activeView === "settings"}
        <div>
          <h1 class="text-xl font-bold">Settings TODO</h1>
        </div>
      {/if}
    </main>
  </div>
</div>
