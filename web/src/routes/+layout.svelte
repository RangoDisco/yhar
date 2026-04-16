<script lang="ts">
    import "./layout.css";
    import favicon from "$lib/assets/favicon.svg";
    import {Toaster} from "$lib/components/ui/sonner";

    let {children} = $props();
</script>

<svelte:head>
    <title>Yhar</title>
    <meta name="description" content="Listening habits stats"/>
    <link rel="icon" href={favicon}/>
</svelte:head>
<section class="min-h-screen md:max-w-300 mx-auto p-6">
    <Toaster/>
    {@render children()}
</section>
