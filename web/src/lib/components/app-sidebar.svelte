<script lang="ts">
    import * as Sidebar from "#lib/components/ui/sidebar/index.js";
    import type {ComponentProps} from "svelte";
    import ModeToggler from "#lib/components/ModeToggler.svelte";
    import Logout from "@lucide/svelte/icons/log-out";
    import House from "@lucide/svelte/icons/house";
    import ContentSearch from "#lib/components/ContentSearch.svelte";

    type NavProps = {
        user: {
            id: string;
            username: string;
            role: "USER" | "ADMIN";
            expiresAt: number;
        };
    } & ComponentProps<typeof Sidebar.Root>;

    let {ref = $bindable(null), user, ...restProps}: NavProps = $props();
</script>
<Sidebar.Root
        bind:ref
        collapsible="icon" class="overflow-hidden border-e" {...restProps}
>
    <Sidebar.Header>
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <Sidebar.MenuButton size="lg"
                                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                    {#snippet child({props})}
                        <a href="/" {...props}>
                            <div
                                    class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
                            >
                                Y
                            </div>
                            <div class="grid flex-1 text-start text-sm leading-tight">
                                <span class="truncate font-medium">Yhar</span>
                            </div>
                        </a>
                    {/snippet}
                </Sidebar.MenuButton>
            </Sidebar.MenuItem>
        </Sidebar.Menu>
    </Sidebar.Header>

    <Sidebar.Content>
        <Sidebar.Group>
            <Sidebar.GroupContent>
                <Sidebar.Menu>
                    <Sidebar.MenuItem>
                        <Sidebar.MenuButton tooltipContent="Home"
                        >
                            <a href="/">
                                <House/>
                            </a>
                            <span>Home</span>
                        </Sidebar.MenuButton>
                    </Sidebar.MenuItem>

                    <Sidebar.MenuItem>
                        <ContentSearch/>
                    </Sidebar.MenuItem>
                </Sidebar.Menu>
            </Sidebar.GroupContent>
        </Sidebar.Group>
    </Sidebar.Content>

    <Sidebar.Footer class="flex flex-col items-center">
        <Sidebar.Menu>
            <Sidebar.MenuItem>
                <ModeToggler/>
            </Sidebar.MenuItem>

            <Sidebar.MenuItem>
                {#if user}
                    <form method="POST" action="/auth/logout">
                        <Sidebar.MenuButton>
                            {#snippet child({props})}
                                <button type="submit" {...props}>
                                    <Logout/>
                                    <span>Logout</span>
                                </button>
                            {/snippet}
                        </Sidebar.MenuButton>
                    </form>
                {/if}
            </Sidebar.MenuItem>
        </Sidebar.Menu>
    </Sidebar.Footer>
</Sidebar.Root>
