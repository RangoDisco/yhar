<script lang="ts">

    import {Button} from "$lib/components/ui/button";
    import ModeToggler from "$lib/components/ModeToggler.svelte";

    type Props = {
        user: {
            id: string
            username: string
            role: "USER" | "ADMIN"
            expiresAt: number
        }
    }

    let {user}: Props = $props();

</script>

<nav class="w-full">
    <div class="flex w-full flex-wrap items-center justify-between gap-2">
        <Button href="/" variant="ghost" size="sm">Home</Button>
        <div class="flex items-center gap-2">
            {#if user}
                <div>
                    <form method="POST" action="/auth/logout">
                        <Button size="sm" type="submit">Logout</Button>
                    </form>
                </div>
            {/if}
            <div>
                <ModeToggler/>
            </div>
        </div>
    </div>
</nav>
