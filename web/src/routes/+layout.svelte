<script lang="ts">
    import "./layout.css";
    import favicon from "$lib/assets/favicon.svg";
    import {Toaster} from "$lib/components/ui/sonner";

    import {ModeWatcher} from "mode-watcher";

    let {children} = $props();
</script>

<section class="min-h-screen md:max-w-300 mx-auto p-6">
    <ModeWatcher/>
    <Toaster/>
    {@render children()}
</section>
