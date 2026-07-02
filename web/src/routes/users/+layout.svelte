<script lang="ts">
    import {Button} from "$lib/components/ui/button";
    import ModeToggler from "$lib/components/ModeToggler.svelte";
    import Nav from "$lib/components/Nav.svelte";

    let {children, data} = $props();
</script>

<main class="flex flex-col">
    <div class="flex items-center gap-2 px-6 h-17">
        <Nav user={data.user}/>
    </div>
    <hr/>
    <div class="md:max-w-300 w-full mx-auto p-6">
        {@render children()}
    </div>
</main>
