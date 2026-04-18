<script lang="ts">
    import {Button} from "$lib/components/ui/button";
    import {ImageUp} from "@lucide/svelte";
    import disc from "$lib/assets/placeholders/disc-3.svg";
    import user from "$lib/assets/placeholders/user.jpg";
    import {toast} from "svelte-sonner";

    type Props = {
        pictureUrl: string;
        alt: string;
        contentID: number;
        contentType: "artists" | "albums";
        uploadEnabled?: boolean
    }

    let {pictureUrl, alt, contentID, contentType, uploadEnabled = $bindable(false)}: Props = $props();
    let input: HTMLInputElement | null = $state(null);
    let image: Element | null = $state(null);
    let fallback = $derived(contentType === "artists" ? user : disc);

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
            body: formData,
        });
        if (!res.ok) {
            throw new Error("Unable to upload image!");
        }
    };

</script>

<div class="relative">
    <img class="peer rounded-full aspect-square h-24 {contentType === 'artists' ? 'rounded-full' : 'rounded-md'}"
         src={pictureUrl ?? fallback} bind:this={image} alt={alt} onerror={() => pictureUrl = fallback}/>

    {#if uploadEnabled}
        <Button variant="outline" size="icon"
                class="invisible peer-hover:visible hover:visible absolute right-0 bottom-0 cursor-pointer"
                onclick={handleClick}>
            <ImageUp/>
        </Button>

        <input id="img-file-input" bind:this={input} type="file" onchange={handleChange}
               accept="image/png, image/jpeg, image/gif, image/webp"
               class="hidden"/>
    {/if}

</div>
