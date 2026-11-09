<script lang="ts">
    import {Button} from "#lib/components/ui/button/index.ts";
    import ModeToggler from "#lib/components/ModeToggler.svelte";
    import House from "@lucide/svelte/icons/house";
    import Logout from "@lucide/svelte/icons/log-out";
    import ContentSearch from "#lib/components/ContentSearch.svelte";

    type Props = {
        user: {
            id: string;
            username: string;
            role: "USER" | "ADMIN";
            expiresAt: number;
        };
    };

    let {user}: Props = $props();
</script>

<nav class="bg-background z-1 h-fit w-full md:h-screen md:w-fit sticky left-0 top-0 p-4 md:py-6 border-r">
    <div class="flex h-full md:flex-col flex-wrap items-center justify-between gap-2">
        <div class="flex flex-col gap-2">
            <Button href="/" variant="ghost" size="icon">
                <House/>
            </Button>
            <ContentSearch/>
        </div>
        <div class="flex md:flex-col items-center gap-2">
            <div>
                <ModeToggler/>
            </div>
            {#if user}
                <div>
                    <form method="POST" action="/auth/logout">
                        <Button size="icon" type="submit" variant="destructive">
                            <Logout/>
                        </Button>
                    </form>
                </div>
            {/if}
        </div>
    </div>
</nav>
