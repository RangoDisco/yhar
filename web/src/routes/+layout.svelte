<script lang="ts">
    import "./layout.css";
    import favicon from "$lib/assets/favicon.svg";

    let {children} = $props();
</script>

<svelte:head>
    <title>Yhar</title>
    <meta name="description" content="Listening habits stats"/>
    <link rel="icon" href={favicon}/>
</svelte:head>
<section class="mx-auto p-6">
    {@render children()}
</section>
