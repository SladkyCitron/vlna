<script lang="ts">
  import { activeView, type View } from "$lib/stores/nav";
  import { Compass, Heart, Settings } from "@lucide/svelte";

  const navItems: { id: View; label: string; icon: any }[] = [
    { id: "explore", label: "Explore", icon: Compass },
    { id: "favorites", label: "Favorites", icon: Heart },
  ];

  function setActiveView(view: View) {
    activeView.set(view);
  }
</script>

<aside
  class="w-60 border-r border-border p-3 flex flex-col justify-between shrink-0"
>
  <nav class="space-y-1">
    <div
      class="py-2 px-4 text-sm font-semibold text-muted-foreground tracking-wider"
    >
      Library
    </div>
    {#each navItems as item}
      <button
        type="button"
        class="w-full flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground focus:outline-none"
        class:bg-accent={item.id === $activeView}
        class:text-accent-foreground={item.id === $activeView}
        onclick={() => setActiveView(item.id)}
      >
        <svelte:component this={item.icon} class="w-4 h-4" />
        {item.label}
      </button>
    {/each}
    <div class="border-t border-border mt-2 py-2">
      <button
        type="button"
        class="w-full flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground focus:outline-none"
        class:bg-accent={"settings" === $activeView}
        class:text-accent-foreground={"settings" === $activeView}
        onclick={() => setActiveView("settings")}
      >
        <Settings class="w-4 h-4" />
        Settings
      </button>
    </div>
  </nav>
</aside>
