<script lang="ts">
  import { Window } from "@wailsio/runtime";
  import { Button } from "$lib/components/ui/button";
  import { Minus, Maximize, X } from "@lucide/svelte";
</script>

<header
  class="col-span-2 h-8 wails-drag border-b border-border/40 flex items-center justify-between pl-4"
>
  <div class="flex items-center gap-2 text-sm font-bold">Vlna</div>
  <div class="wails-no-drag flex">
    <Button variant="ghost" size="icon" onclick={Window.Minimise}>
      <Minus />
    </Button>
    <Button variant="ghost" size="icon" onclick={Window.Maximise}>
      <Maximize />
    </Button>
    <Button variant="ghost" size="icon" onclick={Window.Close}>
      <X />
    </Button>
  </div>
</header>
