<script lang="ts">
    import {Button} from "#lib/components/ui/button/index.ts";
    import {ImageUp} from "@lucide/svelte";
    import {toast} from "svelte-sonner";
    import ContentPicture from "#lib/components/ContentPicture.svelte";
    import type {ContentType} from "#lib/types/content.ts";

    type Props = {
        pictureUrl: string;
        alt: string;
        contentID: number;
        contentType: ContentType;
        uploadEnabled?: boolean;
    };

    let {
        pictureUrl,
        alt,
        contentID,
        contentType,
        uploadEnabled = $bindable(false)
    }: Props = $props();

    let input: HTMLInputElement | null = $state(null);
    let image: HTMLImageElement | null = $state(null);

    const handleClick = (e: Event): void => {
        e.preventDefault();
        input?.click();
    };

    const handleChange = async (): Promise<void> => {
        if (!input?.files) {
            return;
        }
        const file = input.files[0];
        if (!file) {
            return;
        }

        const reader = new FileReader();
        reader.addEventListener("load", function () {
            image?.setAttribute("src", reader.result as string);
        });
        reader.readAsDataURL(file);

        try {
            await uploadImage(file);
            toast.success("Image uploaded successfully!");
        } catch (e) {
            image = null;
            toast.error("Unable to upload image");
        }
    };

    const uploadImage = async (file: File): Promise<void> => {
        const formData = new FormData();
        formData.append("image", file);
        formData.append("contentID", String(contentID));
        formData.append("contentType", contentType);

        const res = await fetch("/api/admin/images/upload", {
            method: "POST",
            body: formData
        });
        if (!res.ok) {
            throw new Error("Unable to upload image!");
        }
    };
</script>

<div class="relative">
    <ContentPicture
            class="peer max-w-fit h-24 w-24 md:h-56 md:w-56"
            url={pictureUrl}
            contentType={contentType}
            bind:img={image}
            {alt}
    />

    {#if uploadEnabled}
        <Button
                variant="outline"
                size="icon"
                class="invisible absolute right-0 bottom-0 cursor-pointer peer-hover:visible hover:visible"
                onclick={handleClick}
        >
            <ImageUp/>
        </Button>

        <input
                id="img-file-input"
                bind:this={input}
                type="file"
                onchange={handleChange}
                accept="image/png, image/jpeg, image/gif, image/webp"
                class="hidden"
        />
    {/if}
</div>
