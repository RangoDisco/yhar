<script lang="ts">
    import {Button} from "$lib/components/ui/button";
    import ModeToggler from "$lib/components/ModeToggler.svelte";

    let {children, data} = $props();
</script>

<main class="flex flex-col gap-4">
    <div class="flex flex-col gap-2">
        <nav class="w-full">
            <ul class="flex w-full flex-wrap items-center justify-end gap-2">
                {#if data.user}
                    <li>
                        <form method="POST" action="/auth/logout">
                            <Button size="sm" type="submit">Logout</Button>
                        </form>
                    </li>
                {/if}
                <li>
                    <ModeToggler/>
                </li>
            </ul>
        </nav>
    </div>
    {@render children()}
</main>
