<script lang="ts">
  import TitleBar from "$lib/components/TitleBar.svelte";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { activeView } from "$lib/stores/nav";
</script>

<div
  class="dark bg-background text-foreground grid h-screen w-screen grid-cols-[240px_1fr] grid-rows-[auto_1fr_auto] overflow-hidden select-none"
>
  <TitleBar />

  <div class="col-span-2 flex min-h-0 overflow-hidden">
    <Sidebar />

    <main class="flex-1 overflow-y-auto p-4">
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
