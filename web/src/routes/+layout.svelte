<script lang="ts">
    import "./layout.css";
    import {Toaster} from "$lib/components/ui/sonner";
    import {ModeWatcher} from "mode-watcher";
    import {TooltipProvider} from "$lib/components/ui/tooltip";

    let {children} = $props();
</script>

<section class="min-h-screen">
    <ModeWatcher/>
    <Toaster/>
    <TooltipProvider>
        {@render children()}
    </TooltipProvider>
</section>
