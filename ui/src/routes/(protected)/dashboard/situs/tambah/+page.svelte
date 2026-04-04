<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import "leaflet/dist/leaflet.css";
    import { apiFetch } from "$lib/api";
    import { goto } from "$app/navigation";
    import * as v from "valibot";

    import DetailMasjidForm from "$lib/components/forms/DetailMasjidForm.svelte";
    import DetailPesantrenForm from "$lib/components/forms/DetailPesantrenForm.svelte";
    import DetailMusholaForm from "$lib/components/forms/DetailMusholaForm.svelte";
    import DetailMTForm from "$lib/components/forms/DetailMTForm.svelte";
    import Toast from "$lib/components/ui/Toast.svelte";

    import {
        SitusInputSchema,
        DetailMasjidSchema,
        DetailMusholaSchema,
        DetailMTSchema,
        DetailPesantrenSchema,
        DEFAULT_SITUS_FORM_DATA,
    } from "$lib/schemas/siteInput";
    import { convertToWebP } from "$lib/utils/imageHelpers";
    import { mapFlatErrors, handleEnterKey } from "$lib/utils/formHelpers";

    import type { PageData } from "./$types";
    let { data }: { data: PageData } = $props();
    let jenisSitusList = $derived(data.jenisSitusList);

    let formData = $state<Record<string, any>>(
        structuredClone(DEFAULT_SITUS_FORM_DATA),
    );

    let validationErrors = $state<Record<string, any>>({});
    let currentPhoneInput = $state("");

    const tipologiOptions: Record<string, string[]> = {
        Masjid: [
            "Masjid Besar",
            "Masjid Jami",
            "Masjid Bersejarah",
            "Masjid Publik",
        ],
        Mushola: ["Mushola"],
        Ponpes: [
            "Salafiyah",
            "Kholafiyah",
            "Salafiyah Wajar Dikdas",
            "Muadalah",
        ],
        MT: ["Majelis Taklim"],
        Madrasah: [
            "Madrasah Diniyah",
            "Madrasah Ibtidaiyah",
            "Madrasah Tsanawiyah",
            "Madrasah Aliyah",
        ],
        Lainnya: ["Umum", "Khusus"],
    };

    let selectedJenisSitusName = $state("");
    let mapContainer: HTMLElement;
    let map: any;
    let marker: any;

    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error">("success");

    let selectedImages = $state<
        { file: File; preview: string; webpBlob: Blob }[]
    >([]);
    let isConverting = $state(false);
    let fileInputRef: HTMLInputElement;
    let isSubmitting = $state(false);

    let isMasjid = $derived(
        selectedJenisSitusName.toLowerCase().includes("masjid"),
    );
    let isPesantren = $derived(
        selectedJenisSitusName.toLowerCase().includes("ponpes") ||
            selectedJenisSitusName.toLowerCase().includes("pesantren"),
    );
    let isMusholla = $derived(
        selectedJenisSitusName.toLowerCase().includes("mushola") ||
            selectedJenisSitusName.toLowerCase().includes("musholla"),
    );
    let isMT = $derived(
        selectedJenisSitusName.toLowerCase() === "mt" ||
            selectedJenisSitusName.toLowerCase().includes("majelis"),
    );

    function getTipologiOptions() {
        if (!selectedJenisSitusName) return [];
        for (const [key, options] of Object.entries(tipologiOptions)) {
            if (
                selectedJenisSitusName.toLowerCase().includes(key.toLowerCase())
            )
                return options;
        }
        return ["Umum", "Khusus"];
    }

    onMount(async () => {
        if (mapContainer && !map) {
            const L = await import("leaflet");
            await import("leaflet/dist/leaflet.css");

            delete (L.Icon.Default.prototype as any)._getIconUrl;
            L.Icon.Default.mergeOptions({
                iconRetinaUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
                iconUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
                shadowUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
            });

            map = L.map(mapContainer).setView([-7.1706, 106.5517], 13);
            L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                attribution: "© OpenStreetMap contributors",
            }).addTo(map);
            marker = L.marker([-7.1706, 106.5517], { draggable: true }).addTo(
                map,
            );

            marker.on("dragend", (e: any) => {
                const latlng = e.target.getLatLng();
                formData.latitude = latlng.lat;
                formData.longitude = latlng.lng;
                validationErrors.latitude = null;
            });

            map.on("click", (e: any) => {
                formData.latitude = e.latlng.lat;
                formData.longitude = e.latlng.lng;
                marker.setLatLng(e.latlng);
                validationErrors.latitude = null;
            });
        }
    });

    onDestroy(() => {
        selectedImages.forEach((img) => {
            if (img?.preview) URL.revokeObjectURL(img.preview);
        });
    });

    function handleJenisSitusChange(event: Event) {
        const target = event.target as HTMLSelectElement;
        const selectedId = target.value;
        formData.jenis_situs_id = selectedId;
        validationErrors.jenis_situs_id = null;

        const selected = jenisSitusList.find(
            (item: any) => item.id === selectedId,
        );
        if (selected) {
            selectedJenisSitusName = selected.nama || "";
            formData.jenis_tipologi = "";
            formData.detail = {};

            const nameLower = selectedJenisSitusName.toLowerCase();
            if (nameLower.includes("masjid")) {
                formData.jenis_tipologi = tipologiOptions.Masjid[0];
            } else if (
                nameLower.includes("ponpes") ||
                nameLower.includes("pesantren")
            ) {
                formData.jenis_tipologi = tipologiOptions.Ponpes[0];
            } else if (
                nameLower.includes("mushola") ||
                nameLower.includes("musholla")
            ) {
                formData.jenis_tipologi = tipologiOptions.Mushola[0];
            } else if (nameLower === "mt" || nameLower.includes("majelis")) {
                formData.jenis_tipologi = tipologiOptions.MT[0];
            }
        }
    }

    function handlePhoneInput(event: Event) {
        const target = event.target as HTMLInputElement;
        target.value = target.value.replace(/[^0-9]/g, "");
        currentPhoneInput = target.value;
    }

    function handleInstansiPhoneInput(event: Event) {
        const target = event.target as HTMLInputElement;
        target.value = target.value.replace(/[^0-9]/g, "");
        formData.nomor_telepon = target.value;
        validationErrors.nomor_telepon = null;
    }

    function addPhone() {
        const phone = currentPhoneInput.trim();
        if (phone && !formData.nomor_telpon_pengurus.includes(phone)) {
            formData.nomor_telpon_pengurus = [
                ...formData.nomor_telpon_pengurus,
                phone,
            ];
            currentPhoneInput = "";
            validationErrors.nomor_telpon_pengurus = null;
        }
    }

    function removePhone(index: number) {
        formData.nomor_telpon_pengurus = formData.nomor_telpon_pengurus.filter(
            (_: any, i: number) => i !== index,
        );
    }

    async function getCurrentLocation(event: Event) {
        if (!navigator.geolocation) {
            toastMessage = "Geolocation tidak didukung";
            toastType = "error";
            showToast = true;
            return;
        }
        const btn = event.currentTarget as HTMLButtonElement;
        const originalText = btn.innerHTML;
        btn.innerHTML =
            '<span class="loading loading-spinner loading-sm text-primary"></span> Mencari...';
        btn.disabled = true;

        try {
            const position = await new Promise<GeolocationPosition>(
                (resolve, reject) => {
                    navigator.geolocation.getCurrentPosition(resolve, reject, {
                        enableHighAccuracy: true,
                        timeout: 15000,
                        maximumAge: 0,
                    });
                },
            );
            formData.latitude = position.coords.latitude;
            formData.longitude = position.coords.longitude;
            validationErrors.latitude = null;

            if (map && marker) {
                marker.setLatLng([formData.latitude, formData.longitude]);
                map.flyTo([formData.latitude, formData.longitude], 16, {
                    animate: true,
                    duration: 1.5,
                });
            }
        } catch (error) {
            toastMessage = "Gagal mendapatkan lokasi";
            toastType = "error";
            showToast = true;
        } finally {
            btn.innerHTML = originalText;
            btn.disabled = false;
        }
    }

    async function handleImageSelect(event: Event) {
        const target = event.target as HTMLInputElement;
        const files = Array.from(target.files || []);
        if (files.length === 0) return;

        isConverting = true;
        try {
            for (const file of files) {
                if (!file.type.startsWith("image/")) continue;
                const preview = URL.createObjectURL(file);
                const webpBlob = await convertToWebP(file);
                selectedImages.push({ file, preview, webpBlob });
            }
        } catch (error) {
            console.error(error);
            toastMessage = "Gagal memproses gambar";
            toastType = "error";
            showToast = true;
        } finally {
            isConverting = false;
        }
    }

    function removeImage(index: number) {
        if (selectedImages[index]?.preview)
            URL.revokeObjectURL(selectedImages[index].preview);
        selectedImages.splice(index, 1);
        if (selectedImages.length === 0 && fileInputRef)
            fileInputRef.value = "";
    }

    async function uploadImages(situsId: string) {
        if (selectedImages.length === 0) return;
        const formDataImages = new FormData();
        for (const img of selectedImages) {
            formDataImages.append(
                "images",
                img.webpBlob,
                `image-${Date.now()}.webp`,
            );
        }
        await apiFetch(`/api/situs/${situsId}/foto`, {
            method: "POST",
            body: formDataImages,
        });
    }

    async function handleSubmit() {
        let errors: Record<string, any> = {};
        let hasError = false;

        const mainResult = v.safeParse(SitusInputSchema, formData);
        if (!mainResult.success) {
            hasError = true;
            mapFlatErrors(v.flatten(mainResult.issues), errors);
        }

        let detailResult;
        if (isMasjid)
            detailResult = v.safeParse(DetailMasjidSchema, formData.detail);
        else if (isMusholla)
            detailResult = v.safeParse(DetailMusholaSchema, formData.detail);
        else if (isMT)
            detailResult = v.safeParse(DetailMTSchema, formData.detail);
        else if (isPesantren)
            detailResult = v.safeParse(DetailPesantrenSchema, formData.detail);

        if (detailResult && !detailResult.success) {
            hasError = true;
            mapFlatErrors(v.flatten(detailResult.issues), errors);
        }

        if (hasError) {
            validationErrors = errors;
            toastMessage = "Tolong periksa kembali form isian Anda";
            toastType = "error";
            showToast = true;
            setTimeout(() => (showToast = false), 3000);
            window.scrollTo({ top: 0, behavior: "smooth" });
            return;
        }

        validationErrors = {};
        isSubmitting = true;

        try {
            const response = await apiFetch("/api/situs", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "Gagal menyimpan data ke server");
            }

            const data = await response.json();

            const situsId = data.id || data.data?.id;
            if (situsId && selectedImages.length > 0) {
                await uploadImages(situsId);
            }

            toastMessage = "Data situs berhasil disimpan";
            toastType = "success";
            showToast = true;

            setTimeout(() => {
                goto("/dashboard/situs");
            }, 2000);
        } catch (error: any) {
            console.error("Failed to submit form:", error);
            toastMessage =
                error.message || "Terjadi kesalahan. Silakan coba lagi.";
            toastType = "error";
            showToast = true;
            setTimeout(() => (showToast = false), 5000);
        } finally {
            isSubmitting = false;
        }
    }
</script>

<svelte:head>
    <title>Tambah Situs Keagamaan | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
</svelte:head>

<div class="min-h-screen bg-base-200/50 p-4 sm:p-8">
    <div class="mx-auto max-w-4xl">
        <div class="mb-8 flex items-center gap-3">
            <div
                class="w-1.5 h-10 rounded-full bg-linear-to-b from-primary to-secondary"
            ></div>
            <div>
                <h1
                    class="text-3xl font-extrabold tracking-tight text-base-content"
                >
                    Pendataan Situs Keagamaan
                </h1>
                <p class="mt-1 text-sm font-medium text-base-content/60">
                    Formulir pendaftaran situs keagamaan baru
                </p>
            </div>
        </div>

        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleSubmit();
            }}
            class="space-y-6"
        >
            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
            >
                <div class="card-body p-5 sm:p-6">
                    <div class="flex items-center gap-3 mb-5">
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-primary/20 to-secondary/20 text-primary font-bold shadow-inner"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                                />
                            </svg>
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Kategori Situs
                        </h2>
                    </div>

                    <div class="form-control w-full">
                        <label class="label" for="jenis_situs_id">
                            <span
                                class="label-text font-bold text-base-content/80"
                                >Pilih Jenis Situs <span class="text-error"
                                    >*</span
                                ></span
                            >
                        </label>
                        <select
                            id="jenis_situs_id"
                            class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.jenis_situs_id
                                ? 'select-error'
                                : ''}"
                            bind:value={formData.jenis_situs_id}
                            onchange={handleJenisSitusChange}
                        >
                            <option value="" disabled
                                >Pilih jenis situs...</option
                            >
                            {#each jenisSitusList as jenis (jenis.id)}
                                <option value={jenis.id}>{jenis.nama}</option>
                            {/each}
                        </select>
                        {#if validationErrors.jenis_situs_id}
                            <span class="mt-1.5 text-xs font-medium text-error"
                                >{validationErrors.jenis_situs_id}</span
                            >
                        {/if}
                    </div>
                </div>
            </div>

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
            >
                <div class="card-body p-5 sm:p-6">
                    <div
                        class="flex items-center gap-3 mb-5 border-b border-base-200 pb-4"
                    >
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-info/20 to-primary/20 text-info font-bold shadow-inner"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                />
                            </svg>
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Informasi Umum
                        </h2>
                    </div>

                    <div class="grid grid-cols-1 gap-5">
                        <div class="form-control">
                            <label class="label pt-0" for="nama">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Nama Situs <span class="text-error">*</span
                                    ></span
                                >
                            </label>
                            <input
                                id="nama"
                                type="text"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nama
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.nama}
                                oninput={() => (validationErrors.nama = null)}
                            />
                            {#if validationErrors.nama}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.nama}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="nomor_pengurus">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >No. Telepon Pengurus <span
                                        class="text-error">*</span
                                    ></span
                                >
                            </label>
                            <div class="flex gap-3">
                                <input
                                    id="nomor_pengurus"
                                    type="tel"
                                    placeholder="Contoh: 0812345678"
                                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nomor_telpon_pengurus
                                        ? 'input-error'
                                        : ''}"
                                    value={currentPhoneInput}
                                    oninput={handlePhoneInput}
                                    onkeydown={(e) => {
                                        handleEnterKey(e, addPhone);
                                    }}
                                />
                                <button
                                    type="button"
                                    class="btn btn-primary rounded-xl px-6 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-12"
                                    onclick={addPhone}
                                >
                                    Tambah
                                </button>
                            </div>
                            <div class="mt-3 flex flex-wrap gap-2">
                                {#each formData.nomor_telpon_pengurus as phone, index}
                                    <div
                                        class="badge font-bold gap-2 px-4 py-3.5 rounded-xl bg-primary/10 text-primary border-primary/20"
                                    >
                                        {phone}
                                        <button
                                            type="button"
                                            title="Hapus nomor"
                                            onclick={() => removePhone(index)}
                                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 text-primary/60 hover:bg-primary/20 hover:text-primary transition-colors"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                class="h-3 w-3"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke="currentColor"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12"
                                                />
                                            </svg>
                                        </button>
                                    </div>
                                {/each}
                            </div>
                            {#if validationErrors.nomor_telpon_pengurus}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.nomor_telpon_pengurus}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="nomor_telepon">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Nomor Telepon Instansi/Situs <span
                                        class="text-xs font-normal text-base-content/50 ml-1"
                                        >(Opsional)</span
                                    ></span
                                >
                            </label>
                            <input
                                id="nomor_telepon"
                                type="tel"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nomor_telepon
                                    ? 'input-error'
                                    : ''}"
                                value={formData.nomor_telepon}
                                oninput={handleInstansiPhoneInput}
                            />
                            {#if validationErrors.nomor_telepon}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.nomor_telepon}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="jenis_tipologi">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Jenis Tipologi <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            {#if getTipologiOptions().length <= 1}
                                <input
                                    id="jenis_tipologi"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed text-base-content/70"
                                    value={formData.jenis_tipologi}
                                    readonly
                                />
                            {:else}
                                <select
                                    id="jenis_tipologi"
                                    class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.jenis_tipologi
                                        ? 'select-error'
                                        : ''}"
                                    bind:value={formData.jenis_tipologi}
                                    onchange={() =>
                                        (validationErrors.jenis_tipologi =
                                            null)}
                                >
                                    {#each getTipologiOptions() as option}
                                        <option value={option}>{option}</option>
                                    {/each}
                                </select>
                            {/if}
                        </div>

                        <div
                            class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-5"
                        >
                            <div class="form-control">
                                <label class="label" for="email">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Email <span
                                            class="text-xs font-normal text-base-content/50 ml-1"
                                            >(Opsional)</span
                                        ></span
                                    >
                                </label>
                                <input
                                    id="email"
                                    type="email"
                                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.email
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={formData.email}
                                    oninput={() =>
                                        (validationErrors.email = null)}
                                />
                                {#if validationErrors.email}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.email}</span
                                    >
                                {/if}
                            </div>
                            <div class="form-control">
                                <label class="label" for="website">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Website <span
                                            class="text-xs font-normal text-base-content/50 ml-1"
                                            >(Opsional)</span
                                        ></span
                                    >
                                </label>
                                <input
                                    id="website"
                                    type="url"
                                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.website
                                        ? 'input-error'
                                        : ''}"
                                    placeholder="https://..."
                                    bind:value={formData.website}
                                    oninput={() =>
                                        (validationErrors.website = null)}
                                />
                                {#if validationErrors.website}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.website}</span
                                    >
                                {/if}
                            </div>
                        </div>

                        <div
                            class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-5"
                        >
                            <div class="form-control">
                                <label class="label" for="nomor_badan_hukum">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Nomor Badan Hukum <span
                                            class="text-xs font-normal text-base-content/50 ml-1"
                                            >(Opsional)</span
                                        ></span
                                    >
                                </label>
                                <input
                                    id="nomor_badan_hukum"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nomor_badan_hukum
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={formData.nomor_badan_hukum}
                                    oninput={() =>
                                        (validationErrors.nomor_badan_hukum =
                                            null)}
                                />
                                {#if validationErrors.nomor_badan_hukum}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.nomor_badan_hukum}</span
                                    >
                                {/if}
                            </div>
                            <div class="form-control">
                                <label class="label" for="tahun_berdiri">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Tahun Berdiri <span class="text-error"
                                            >*</span
                                        ></span
                                    >
                                </label>
                                <input
                                    id="tahun_berdiri"
                                    type="number"
                                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.tahun_berdiri
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={formData.tahun_berdiri}
                                    oninput={() =>
                                        (validationErrors.tahun_berdiri = null)}
                                />
                                {#if validationErrors.tahun_berdiri}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.tahun_berdiri}</span
                                    >
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
            >
                <div class="card-body p-5 sm:p-6">
                    <div
                        class="flex items-center gap-3 mb-5 border-b border-base-200 pb-4"
                    >
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-success/20 to-info/20 text-success font-bold shadow-inner"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                />
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                />
                            </svg>
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Lokasi & Peta
                        </h2>
                    </div>

                    <div class="flex flex-col gap-5">
                        <div
                            class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-5"
                        >
                            <div class="form-control">
                                <label class="label pt-0" for="provinsi">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Provinsi</span
                                    >
                                </label>
                                <input
                                    id="provinsi"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed text-base-content/70"
                                    value={formData.provinsi}
                                    readonly
                                />
                            </div>
                            <div class="form-control">
                                <label class="label pt-0" for="kabupaten_kota">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Kabupaten/Kota</span
                                    >
                                </label>
                                <input
                                    id="kabupaten_kota"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed text-base-content/70"
                                    value={formData.kabupaten_kota}
                                    readonly
                                />
                            </div>
                            <div class="form-control">
                                <label class="label" for="kecamatan">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Kecamatan</span
                                    >
                                </label>
                                <input
                                    id="kecamatan"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed text-base-content/70"
                                    value={formData.kecamatan}
                                    readonly
                                />
                            </div>
                            <div class="form-control">
                                <label class="label" for="desa">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Desa/Kelurahan <span class="text-error"
                                            >*</span
                                        ></span
                                    >
                                </label>
                                <input
                                    id="desa"
                                    type="text"
                                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.desa
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={formData.desa}
                                    oninput={() =>
                                        (validationErrors.desa = null)}
                                />
                                {#if validationErrors.desa}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.desa}</span
                                    >
                                {/if}
                            </div>
                        </div>

                        <div class="form-control">
                            <label class="label" for="alamat_lengkap">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Alamat Lengkap <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            <textarea
                                id="alamat_lengkap"
                                class="textarea textarea-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-24 py-3 {validationErrors.alamat_lengkap
                                    ? 'textarea-error'
                                    : ''}"
                                bind:value={formData.alamat_lengkap}
                                oninput={() =>
                                    (validationErrors.alamat_lengkap = null)}
                                rows="3"
                            ></textarea>
                            {#if validationErrors.alamat_lengkap}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.alamat_lengkap}</span
                                >
                            {/if}
                        </div>

                        <div class="mt-4">
                            <button
                                type="button"
                                class="btn btn-secondary rounded-xl px-6 shadow-lg shadow-secondary/20 transition-transform hover:-translate-y-0.5 min-h-12 mb-5 w-full sm:w-auto"
                                onclick={getCurrentLocation}
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-5 w-5 mr-1.5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                    />
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                    />
                                </svg>
                                Deteksi Lokasi Saat Ini
                            </button>

                            <div
                                class="grid grid-cols-1 gap-4 md:grid-cols-2 mb-4"
                            >
                                <div class="form-control">
                                    <label class="label py-2" for="latitude">
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Latitude <span class="text-error"
                                                >*</span
                                            ></span
                                        >
                                    </label>
                                    <input
                                        required
                                        id="latitude"
                                        type="number"
                                        step="any"
                                        class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed {validationErrors.latitude
                                            ? 'input-error'
                                            : ''}"
                                        value={formData.latitude ?? ""}
                                        readonly
                                    />
                                </div>
                                <div class="form-control">
                                    <label class="label py-2" for="longitude">
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Longitude <span class="text-error"
                                                >*</span
                                            ></span
                                        >
                                    </label>
                                    <input
                                        required
                                        id="longitude"
                                        type="number"
                                        step="any"
                                        class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed {validationErrors.longitude
                                            ? 'input-error'
                                            : ''}"
                                        value={formData.longitude ?? ""}
                                        readonly
                                    />
                                </div>
                            </div>

                            {#if validationErrors.latitude || validationErrors.longitude}
                                <span
                                    class="mb-4 block text-xs font-medium text-error"
                                >
                                    Silakan tentukan titik lokasi pada peta
                                </span>
                            {/if}

                            <div
                                class="h-72 w-full overflow-hidden rounded-2xl border border-base-200 shadow-inner bg-base-200/30 sm:h-80 md:h-96 relative z-0"
                            >
                                <div
                                    bind:this={mapContainer}
                                    class="h-full w-full"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
            >
                <div class="card-body p-5 sm:p-6">
                    <div
                        class="flex items-center gap-3 mb-5 border-b border-base-200 pb-4"
                    >
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-warning/20 to-error/20 text-warning font-bold shadow-inner"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                />
                            </svg>
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Kapasitas & Legalitas
                        </h2>
                    </div>

                    <div class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-5">
                        <div class="form-control">
                            <label class="label pt-0" for="luas_tanah">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Luas Tanah (m²) <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            <input
                                id="luas_tanah"
                                type="number"
                                min="0"
                                step="any"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.luas_tanah
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.luas_tanah}
                                oninput={() =>
                                    (validationErrors.luas_tanah = null)}
                            />
                            {#if validationErrors.luas_tanah}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.luas_tanah}</span
                                >
                            {/if}
                        </div>
                        <div class="form-control">
                            <label class="label pt-0" for="luas_bangunan">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Luas Bangunan (m²) <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            <input
                                id="luas_bangunan"
                                type="number"
                                min="0"
                                step="any"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.luas_bangunan
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.luas_bangunan}
                                oninput={() =>
                                    (validationErrors.luas_bangunan = null)}
                            />
                            {#if validationErrors.luas_bangunan}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.luas_bangunan}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="status_tanah">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Status Tanah <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            <select
                                id="status_tanah"
                                class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.status_tanah
                                    ? 'select-error'
                                    : ''}"
                                bind:value={formData.status_tanah}
                                onchange={() =>
                                    (validationErrors.status_tanah = null)}
                            >
                                <option value="">Pilih status...</option>
                                <option value="Wakaf">Wakaf</option>
                                <option value="Milik Sendiri"
                                    >Milik Sendiri</option
                                >
                                <option value="Hibah">Hibah</option>
                                <option value="Lainnya">Lainnya</option>
                            </select>
                            {#if validationErrors.status_tanah}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.status_tanah}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="daya_tampung_max">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Daya Tampung Max <span class="text-error"
                                        >*</span
                                    ></span
                                >
                            </label>
                            <input
                                id="daya_tampung_max"
                                type="number"
                                min="0"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.daya_tampung_max
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.daya_tampung_max}
                                oninput={() =>
                                    (validationErrors.daya_tampung_max = null)}
                            />
                            {#if validationErrors.daya_tampung_max}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.daya_tampung_max}</span
                                >
                            {/if}
                        </div>

                        <div class="form-control">
                            <label class="label" for="nomor_aiw">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >Nomor AIW <span
                                        class="text-xs font-normal text-base-content/50 ml-1"
                                        >(Opsional)</span
                                    ></span
                                >
                            </label>
                            <input
                                id="nomor_aiw"
                                type="text"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nomor_aiw
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.nomor_aiw}
                                oninput={() =>
                                    (validationErrors.nomor_aiw = null)}
                            />
                            {#if validationErrors.nomor_aiw}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.nomor_aiw}</span
                                >
                            {/if}
                        </div>
                        <div class="form-control">
                            <label class="label" for="nomor_sertifikat_wakaf">
                                <span
                                    class="label-text font-bold text-base-content/80"
                                    >No. Sertifikat Wakaf <span
                                        class="text-xs font-normal text-base-content/50 ml-1"
                                        >(Opsional)</span
                                    ></span
                                >
                            </label>
                            <input
                                id="nomor_sertifikat_wakaf"
                                type="text"
                                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.nomor_sertifikat_wakaf
                                    ? 'input-error'
                                    : ''}"
                                bind:value={formData.nomor_sertifikat_wakaf}
                                oninput={() =>
                                    (validationErrors.nomor_sertifikat_wakaf =
                                        null)}
                            />
                            {#if validationErrors.nomor_sertifikat_wakaf}
                                <span
                                    class="mt-1.5 text-xs font-medium text-error"
                                    >{validationErrors.nomor_sertifikat_wakaf}</span
                                >
                            {/if}
                        </div>
                    </div>
                </div>
            </div>

            {#if isMasjid}
                <DetailMasjidForm
                    bind:detail={formData.detail}
                    bind:errors={validationErrors}
                />
            {/if}
            {#if isPesantren}
                <DetailPesantrenForm
                    bind:detail={formData.detail}
                    bind:errors={validationErrors}
                />
            {/if}
            {#if isMusholla}
                <DetailMusholaForm
                    bind:detail={formData.detail}
                    bind:errors={validationErrors}
                />
            {/if}
            {#if isMT}
                <DetailMTForm
                    bind:detail={formData.detail}
                    bind:errors={validationErrors}
                />
            {/if}

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
            >
                <div class="card-body p-5 sm:p-6">
                    <div
                        class="flex items-center gap-3 mb-5 border-b border-base-200 pb-4"
                    >
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-secondary/20 to-accent/20 text-secondary font-bold shadow-inner"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                />
                            </svg>
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Foto Situs
                        </h2>
                    </div>

                    <div class="form-control">
                        <label class="label pt-0" for="images">
                            <span
                                class="label-text font-bold text-base-content/80"
                                >Unggah Foto</span
                            >
                        </label>
                        <div class="mb-3 flex flex-col gap-3">
                            <input
                                id="images"
                                type="file"
                                accept="image/*"
                                multiple
                                class="file-input file-input-bordered w-full rounded-xl bg-base-200/30 transition-colors min-h-12"
                                onchange={handleImageSelect}
                                disabled={isConverting}
                                bind:this={fileInputRef}
                            />
                            {#if isConverting}
                                <div
                                    class="flex items-center gap-3 text-sm text-base-content/70 bg-base-200/50 rounded-xl p-4 font-medium border border-base-200"
                                >
                                    <span
                                        class="loading loading-sm loading-spinner text-primary"
                                    ></span>
                                    Mengkonversi gambar ke format WebP untuk menghemat
                                    penyimpanan...
                                </div>
                            {/if}
                        </div>

                        {#if selectedImages.length > 0}
                            <div
                                class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4"
                            >
                                {#each selectedImages as img, index (img.preview)}
                                    <div
                                        class="relative overflow-hidden rounded-xl border border-base-200 shadow-sm group"
                                    >
                                        <img
                                            src={img.preview}
                                            alt="Preview"
                                            class="h-32 w-full object-cover transition-transform duration-500 group-hover:scale-110"
                                        />
                                        <button
                                            title="Hapus foto"
                                            type="button"
                                            class="btn btn-circle btn-xs btn-error absolute top-2 right-2 text-white shadow-lg opacity-90 hover:opacity-100"
                                            onclick={() => removeImage(index)}
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                class="h-3 w-3"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke="currentColor"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12"
                                                />
                                            </svg>
                                        </button>
                                        <div
                                            class="absolute bottom-2 right-2 rounded-lg bg-black/70 backdrop-blur-sm px-2.5 py-1 text-xs font-bold tracking-wider text-white"
                                        >
                                            WebP
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </div>
            </div>

            <div
                class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8 pt-4"
            >
                <button
                    onclick={() => goto("/dashboard/situs")}
                    type="button"
                    class="btn btn-ghost rounded-xl px-8 min-h-12 font-bold text-base-content/60 hover:bg-base-200/50 hover:text-base-content"
                    disabled={isSubmitting}
                >
                    Batal
                </button>
                <button
                    type="submit"
                    class="btn btn-primary rounded-xl px-10 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-12 font-extrabold text-base-100"
                    disabled={isSubmitting}
                >
                    {#if isSubmitting}
                        <span class="loading loading-spinner mr-2"></span> Menyimpan
                        Data...
                    {:else}
                        Simpan Data Situs
                    {/if}
                </button>
            </div>
        </form>
    </div>

    <Toast
        show={showToast}
        message={toastMessage}
        type={toastType}
        onclose={() => (showToast = false)}
    />
</div>

<style>
    :global(.leaflet-container) {
        z-index: 0;
    }
</style>
