<script lang="ts">
    import Nav from "#lib/components/Nav.svelte";

    let {children, data} = $props();
</script>

<Nav user={data.user}/>
<div class="mx-auto w-full p-10 lg:max-w-500">
    {@render children()}
</div>
