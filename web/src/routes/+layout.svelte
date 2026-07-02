<script lang="ts">
    import "./layout.css";
    import {Toaster} from "$lib/components/ui/sonner";
    import {ModeWatcher} from "mode-watcher";

    let {children} = $props();
</script>

<section class="min-h-screen">
    <ModeWatcher/>
    <Toaster/>
    {@render children()}
</section>
