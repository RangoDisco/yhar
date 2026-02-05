<script>
    import * as Avatar from "$lib/components/ui/avatar/index";

    let {
        title,
        subtitle = $bindable(null),
        pictureUrl,
        contentType
    } = $props();
</script>

<article class="flex items-center gap-4 justify-between">
    <div class="flex gap-4 items-center">
    <h3 class="text-lg">1</h3>
    <div class="flex gap-2">
        <Avatar.Root class="{contentType !== 'artists' ? 'rounded-md' : ''} h-12 w-12 ">
            <Avatar.Image src={pictureUrl}
                          alt={`${title}'s picture`}/>
            <Avatar.Fallback>{title}</Avatar.Fallback>
        </Avatar.Root>
        <div class="flex flex-col">
            <p>{title}</p>
            <p class="text-sm text-muted-foreground">{subtitle}</p>
        </div>
    </div>
    </div>
    <p class="text-muted-foreground">92</p>
</article>
