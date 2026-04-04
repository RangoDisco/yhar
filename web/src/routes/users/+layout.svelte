<script lang="ts">
    import {Button} from "$lib/components/ui/button";

    let {children, data} = $props();
</script>

<main class="flex flex-col gap-4">
    {#if data.isLoggedIn}
        <div class="flex flex-col gap-2">
            <nav class="w-full">
                <ul class="flex w-full flex-wrap items-center justify-end gap-2">
                    <li>
                        <form method="POST" action="/auth/logout">
                            <Button size="sm" type="submit">Logout</Button>
                        </form>
                    </li>
                </ul>
            </nav>
        </div>
    {/if}
    {@render children()}
</main>
