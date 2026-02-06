<script>
    import * as Avatar from "$lib/components/ui/avatar/index";
    import {Button} from "$lib/components/ui/button/index";
    import {Trash} from "@lucide/svelte";

    let {
        title,
        album,
        scrobbledAt,
        pictureUrl,
    } = $props();
</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-2 items-center">
        <Avatar.Root class="rounded-md h-8 w-8 ">
            <Avatar.Image src={pictureUrl}
                          alt={`${title}'s picture`}/>
            <Avatar.Fallback>{title}</Avatar.Fallback>
        </Avatar.Root>
        <div class="flex flex-col">
            <div class="flex gap-2">
                <p>{title}</p>
                <p class="text-sm text-muted-foreground">{album}</p>
            </div>
            <p class="text-sm text-muted-foreground">{scrobbledAt}</p>
        </div>
    </div>
    <Button variant="outline" size="icon">
        <!--        TODO: Handle delete-->
        <Trash/>
    </Button>
</article>
