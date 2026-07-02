<script lang="ts">
    import favicon from "$lib/assets/favicon.svg";

    type Props = {
        title: string;
        description?: string | null;
    }

    let {title, description = $bindable()}: Props = $props();
</script>

<svelte:head>
    <title>{title} - Yhar</title>
    <meta name="description" content={description ?? "Listening habits stats"}/>
    <link rel="icon" href={favicon}/>
</svelte:head>
