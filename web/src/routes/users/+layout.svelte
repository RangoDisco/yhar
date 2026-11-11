<script lang="ts">
    import Nav from "#lib/components/Nav.svelte";
    import * as Sidebar from "#lib/components/ui/sidebar/index.ts";
    import AutoCrumbs from "#lib/components/AutoCrumbs.svelte";

    let {children, data} = $props();
</script>

<Sidebar.Provider open={false}>
    <Nav user={data.user}/>
    <Sidebar.Inset>
        <header class="sticky top-0 flex shrink-0 items-center gap-2 border-b bg-background p-4 z-1">
            <Sidebar.Trigger class="-ms-1"/>
            <AutoCrumbs/>
        </header>
        <div class="mx-auto w-full p-6 lg:max-w-500">
            {@render children()}
        </div>
    </Sidebar.Inset>
</Sidebar.Provider>
