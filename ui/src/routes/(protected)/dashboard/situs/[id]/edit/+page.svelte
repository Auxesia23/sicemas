<script lang="ts">
    import { untrack, onDestroy } from "svelte";
    import "leaflet/dist/leaflet.css";
    import { apiFetch } from "$lib/api";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import * as v from "valibot";

    import DetailMasjidForm from "$lib/components/forms/DetailMasjidForm.svelte";
    import DetailPesantrenForm from "$lib/components/forms/DetailPesantrenForm.svelte";
    import DetailMusholaForm from "$lib/components/forms/DetailMusholaForm.svelte";
    import DetailMTForm from "$lib/components/forms/DetailMTForm.svelte";
    import Toast from "$lib/components/ui/Toast.svelte";
    import ConfirmModal from "$lib/components/ui/ConfirmModal.svelte";

    import type { Foto, SitusDetailResponse } from "$lib/schemas/sites";

    import {
        SitusEditSchema,
        DetailMasjidSchema,
        DetailMusholaSchema,
        DetailMTSchema,
        DetailPesantrenSchema,
    } from "$lib/schemas/siteInput";
    import { mapFlatErrors, handleEnterKey } from "$lib/utils/formHelpers";
    import { convertToWebP } from "$lib/utils/imageHelpers";

    import type { PageData } from "./$types";
    let { data }: { data: PageData } = $props();

    let formData = $state<SitusDetailResponse>(
        untrack(() => structuredClone(data.situs)),
    );

    let validationErrors = $state<Record<string, any>>({});
    let currentPhoneInput = $state("");

    let jenisSitusName = $derived(formData.jenis_situs);

    let existingPhotos = $derived<Foto[]>(formData.foto || []);

    let photosToDelete = $state<string[]>([]);
    let newPhotos = $state<{ file: File; preview: string; webpBlob: Blob }[]>(
        [],
    );

    let isConverting = $state(false);
    let fileInputRef: HTMLInputElement | undefined = $state();

    let map: any = $state(null);
    let marker: any = $state(null);

    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error">("success");
    let showConfirmModal = $state(false);
    let isSubmitting = $state(false);
    let isLoading = $state(false);

    let isMasjid = $derived(
        jenisSitusName.toLowerCase().trim().includes("masjid"),
    );
    let isPesantren = $derived(
        jenisSitusName.toLowerCase().trim().includes("ponpes") ||
            jenisSitusName.toLowerCase().trim().includes("pesantren"),
    );
    let isMusholla = $derived(
        jenisSitusName.toLowerCase().trim().includes("mushola") ||
            jenisSitusName.toLowerCase().trim().includes("musholla"),
    );
    let isMT = $derived(
        jenisSitusName.toLowerCase().trim() === "mt" ||
            jenisSitusName.toLowerCase().trim().includes("majelis"),
    );

    const tipologiOptions: Record<string, string[]> = {
        Masjid: [
            "Masjid Besar",
            "Masjid Jami",
            "Masjid Bersejarah",
            "Masjid Publik",
        ],
        Mushola: ["Musholla"],
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

    function getTipologiOptions() {
        if (!jenisSitusName) return [];
        const nameStr = jenisSitusName.toLowerCase().trim();
        if (nameStr === "mt" || nameStr.includes("majelis"))
            return tipologiOptions["MT"];
        if (nameStr.includes("masjid")) return tipologiOptions["Masjid"];
        if (nameStr.includes("mushola") || nameStr.includes("musholla"))
            return tipologiOptions["Mushola"];
        if (nameStr.includes("ponpes") || nameStr.includes("pesantren"))
            return tipologiOptions["Ponpes"];
        if (nameStr.includes("madrasah")) return tipologiOptions["Madrasah"];
        return ["Umum", "Khusus"];
    }

    function mapAction(node: HTMLElement, coords: any) {
        let L: any;
        setTimeout(async () => {
            L = await import("leaflet");
            await import("leaflet/dist/leaflet.css");

            const initialLat = coords.latitude || -6.2088;
            const initialLng = coords.longitude || 106.8456;

            map = L.map(node, {
                dragging: true,
                touchZoom: true,
                scrollWheelZoom: false,
                doubleClickZoom: false,
                boxZoom: false,
                keyboard: false,
                zoomControl: true,
                attributionControl: true,
            }).setView([initialLat, initialLng], 15);

            L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                attribution: "© OpenStreetMap contributors",
            }).addTo(map);

            delete (L.Icon.Default.prototype as any)._getIconUrl;
            L.Icon.Default.mergeOptions({
                iconRetinaUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
                iconUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
                shadowUrl:
                    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
            });

            marker = L.marker([initialLat, initialLng], {
                draggable: true,
            }).addTo(map);

            setTimeout(() => {
                if (map) map.invalidateSize();
            }, 300);

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
        }, 100);

        return {
            update(newCoords: any) {
                if (
                    map &&
                    marker &&
                    newCoords &&
                    newCoords.latitude !== null &&
                    newCoords.longitude !== null
                ) {
                    const newLatLng = new L.LatLng(
                        newCoords.latitude,
                        newCoords.longitude,
                    );
                    if (!marker.getLatLng().equals(newLatLng)) {
                        marker.setLatLng(newLatLng);
                        map.setView(newLatLng, 15);
                    }
                }
            },
            destroy() {
                if (map) {
                    map.remove();
                    map = null;
                    marker = null;
                }
            },
        };
    }

    onDestroy(() => {
        newPhotos.forEach((img) => {
            if (img?.preview) URL.revokeObjectURL(img.preview);
        });
        if (map) {
            map.remove();
            map = null;
        }
    });

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

    async function getCurrentLocation() {
        if (!navigator.geolocation) {
            toastMessage = "Geolocation tidak didukung";
            toastType = "error";
            showToast = true;
            return;
        }
        try {
            const position = await new Promise<GeolocationPosition>(
                (resolve, reject) => {
                    navigator.geolocation.getCurrentPosition(resolve, reject, {
                        enableHighAccuracy: true,
                        timeout: 15000,
                    });
                },
            );
            formData.latitude = position.coords.latitude;
            formData.longitude = position.coords.longitude;
            validationErrors.latitude = null;

            if (map && marker) {
                const latlng = [formData.latitude, formData.longitude];
                marker.setLatLng(latlng);
                map.flyTo(latlng, 16);
            }
        } catch (error) {
            toastMessage = "Gagal mendapatkan lokasi";
            toastType = "error";
            showToast = true;
        }
    }

    async function handlePhotoSelect(event: Event) {
        const target = event.target as HTMLInputElement;
        const files = Array.from(target.files || []);
        if (files.length === 0) return;
        isConverting = true;
        try {
            for (const file of files) {
                if (!file.type.startsWith("image/")) continue;
                const preview = URL.createObjectURL(file);
                const webpBlob = await convertToWebP(file);

                newPhotos = [...newPhotos, { file, preview, webpBlob }];
            }
        } catch (error) {
            toastMessage = "Gagal mengkonversi gambar";
            toastType = "error";
            showToast = true;
        } finally {
            isConverting = false;
            if (fileInputRef) fileInputRef.value = "";
        }
    }

    function removeNewPhoto(index: number) {
        if (newPhotos[index]?.preview)
            URL.revokeObjectURL(newPhotos[index].preview);
        newPhotos = newPhotos.filter((_, i) => i !== index);
    }

    function markPhotoForDeletion(photoId: string) {
        const id = photoId.trim();
        if (!id) return;

        if (!photosToDelete.includes(id)) {
            photosToDelete = [...photosToDelete, id];
        }

        if (formData.foto) {
            formData.foto = formData.foto.filter((p) => p.id !== id);
        }
    }

    function confirmUpdate() {
        let errors: Record<string, any> = {};
        let hasError = false;

        const mainResult = v.safeParse(SitusEditSchema, formData);
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
        showConfirmModal = true;
    }

    async function handleUpdate() {
        isSubmitting = true;
        showConfirmModal = false;

        try {
            const updateResponse = await apiFetch(
                `/api/situs/${page.params.id}`,
                {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(formData),
                },
            );

            if (!updateResponse.ok) {
                const errorText = await updateResponse.text();
                throw new Error(errorText || "Gagal memperbarui data situs");
            }

            const situsId = page.params.id;
            const tasks = [];

            if (newPhotos.length > 0) {
                const fd = new FormData();
                newPhotos.forEach((img) =>
                    fd.append(
                        "images",
                        img.webpBlob,
                        `image-${Date.now()}.webp`,
                    ),
                );
                tasks.push(
                    apiFetch(`/api/situs/${situsId}/foto`, {
                        method: "POST",
                        body: fd,
                    }).then((res) => {
                        if (!res.ok)
                            throw new Error("Gagal mengunggah foto baru");
                    }),
                );
            }

            if (photosToDelete.length > 0) {
                tasks.push(
                    apiFetch(`/api/situs/${situsId}/foto`, {
                        method: "DELETE",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ ids: photosToDelete }),
                    }).then((res) => {
                        if (!res.ok)
                            throw new Error("Gagal menghapus foto lama");
                    }),
                );
            }

            if (tasks.length > 0) await Promise.all(tasks);

            toastMessage = "Data situs berhasil diperbarui";
            toastType = "success";
            showToast = true;

            newPhotos.forEach((img) => {
                if (img?.preview) URL.revokeObjectURL(img.preview);
            });

            setTimeout(() => {
                showToast = false;
                goto("/dashboard/situs");
            }, 2000);
        } catch (error: any) {
            console.error("Update error:", error);
            toastMessage =
                error.message || "Terjadi kesalahan saat memperbarui data";
            toastType = "error";
            showToast = true;
        } finally {
            isSubmitting = false;
        }
    }

    function goBack() {
        goto("/dashboard/situs");
    }
</script>

<svelte:head>
    <title>Edit {formData.nama || "Situs"} | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
</svelte:head>

<div class="relative min-h-screen bg-base-200/50 p-4 sm:p-8">
    <Toast
        show={showToast}
        message={toastMessage}
        type={toastType}
        onclose={() => (showToast = false)}
    />

    <ConfirmModal
        show={showConfirmModal}
        title="Konfirmasi Pembaruan"
        message="Apakah Anda yakin ingin memperbarui data situs ini?"
        confirmText="Simpan Perubahan"
        cancelText="Batal"
        isProcessing={isSubmitting}
        onConfirm={handleUpdate}
        type="warning"
        onCancel={() => (showConfirmModal = false)}
    />

    {#if isLoading}
        <div class="flex h-[80vh] items-center justify-center">
            <span class="loading loading-lg loading-spinner text-primary"
            ></span>
        </div>
    {:else}
        <div class="mx-auto max-w-4xl">
            <div class="mb-8 flex items-start sm:items-center gap-4">
                <button
                    class="btn btn-circle btn-ghost bg-base-200/50 hover:bg-base-300 transition-colors mt-1 sm:mt-0 shadow-sm"
                    onclick={goBack}
                    aria-label="Kembali"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M10 19l-7-7m0 0l7-7m-7 7h18"
                        />
                    </svg>
                </button>
                <div>
                    <div class="flex items-center gap-3 mb-1">
                        <div
                            class="w-1.5 h-8 rounded-full bg-linear-to-b from-primary to-secondary hidden sm:block"
                        ></div>
                        <h1
                            class="text-2xl font-extrabold tracking-tight text-base-content sm:text-3xl leading-tight"
                        >
                            Edit Situs Keagamaan
                        </h1>
                    </div>
                    <p
                        class="text-sm font-medium text-base-content/60 sm:text-base sm:ml-5"
                    >
                        Perbarui informasi situs keagamaan
                    </p>
                </div>
            </div>

            <form
                class="flex flex-col gap-6"
                onsubmit={(e) => {
                    e.preventDefault();
                    confirmUpdate();
                }}
            >
                <div
                    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
                >
                    <div class="card-body p-5 sm:p-6">
                        <div
                            class="flex items-center gap-3 mb-5 border-b border-base-200 pb-4"
                        >
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
                                Jenis Situs
                            </h2>
                        </div>
                        <div class="form-control w-full">
                            <label
                                class="label pt-0 pb-2"
                                for="jenis_situs_display"
                            >
                                <span
                                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                                    >Jenis Situs (Tidak dapat diubah)</span
                                >
                            </label>
                            <input
                                id="jenis_situs_display"
                                type="text"
                                class="input input-bordered w-full rounded-xl bg-base-200/50 font-semibold text-primary min-h-12 cursor-not-allowed"
                                value={jenisSitusName}
                                readonly
                            />
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

                        <div class="flex flex-col gap-5">
                            <div class="form-control">
                                <label class="label pt-0 pb-2" for="nama">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Nama Situs <span class="text-error"
                                            >*</span
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
                                    oninput={() =>
                                        (validationErrors.nama = null)}
                                />
                                {#if validationErrors.nama}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.nama}</span
                                    >
                                {/if}
                            </div>

                            <div class="form-control">
                                <label class="label py-2" for="nomor_pengurus">
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
                                        onkeydown={(e) =>
                                            handleEnterKey(e, addPhone)}
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
                                    {#each formData.nomor_telpon_pengurus || [] as phone, index}
                                        <div
                                            class="badge font-bold gap-2 px-4 py-3.5 rounded-xl bg-primary/10 text-primary border-primary/20"
                                        >
                                            {phone}
                                            <button
                                                type="button"
                                                title="Hapus nomor"
                                                onclick={() =>
                                                    removePhone(index)}
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

                            <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                                <div class="form-control">
                                    <label
                                        class="label py-2"
                                        for="nomor_telepon"
                                    >
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
                                    <label
                                        class="label py-2"
                                        for="jenis_tipologi"
                                    >
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Jenis Tipologi <span
                                                class="text-error">*</span
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
                                            <option value="" disabled
                                                >Pilih jenis tipologi...</option
                                            >
                                            {#each getTipologiOptions() as option}
                                                <option value={option}
                                                    >{option}</option
                                                >
                                            {/each}
                                        </select>
                                    {/if}
                                    {#if validationErrors.jenis_tipologi}
                                        <span
                                            class="mt-1.5 text-xs font-medium text-error"
                                            >{validationErrors.jenis_tipologi}</span
                                        >
                                    {/if}
                                </div>
                            </div>

                            <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                                <div class="form-control">
                                    <label class="label py-2" for="email">
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
                                    <label class="label py-2" for="website">
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

                            <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                                <div class="form-control">
                                    <label
                                        class="label py-2"
                                        for="nomor_badan_hukum"
                                    >
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
                                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12"
                                        bind:value={formData.nomor_badan_hukum}
                                    />
                                </div>
                                <div class="form-control">
                                    <label
                                        class="label py-2"
                                        for="tahun_berdiri"
                                    >
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Tahun Berdiri <span
                                                class="text-error">*</span
                                            ></span
                                        >
                                    </label>
                                    <input
                                        id="tahun_berdiri"
                                        type="number"
                                        min="0"
                                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-12 {validationErrors.tahun_berdiri
                                            ? 'input-error'
                                            : ''}"
                                        bind:value={formData.tahun_berdiri}
                                        oninput={() =>
                                            (validationErrors.tahun_berdiri =
                                                null)}
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
                            <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                                <div class="form-control">
                                    <label
                                        class="label pt-0 pb-2"
                                        for="provinsi"
                                    >
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
                                    <label
                                        class="label pt-0 pb-2"
                                        for="kabupaten_kota"
                                    >
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
                                    <label class="label py-2" for="kecamatan">
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
                                    <label class="label py-2" for="desa">
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Desa/Kelurahan <span
                                                class="text-error">*</span
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
                                <label class="label py-2" for="alamat_lengkap">
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
                                    rows="3"
                                    oninput={() =>
                                        (validationErrors.alamat_lengkap =
                                            null)}
                                ></textarea>
                                {#if validationErrors.alamat_lengkap}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.alamat_lengkap}</span
                                    >
                                {/if}
                            </div>

                            <div class="mt-2">
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
                                        <label
                                            class="label py-2"
                                            for="latitude"
                                        >
                                            <span
                                                class="label-text font-bold text-base-content/80"
                                                >Latitude <span
                                                    class="text-error">*</span
                                                ></span
                                            >
                                        </label>
                                        <input
                                            id="latitude"
                                            type="number"
                                            step="any"
                                            class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed {validationErrors.latitude
                                                ? 'input-error'
                                                : ''}"
                                            bind:value={formData.latitude}
                                            readonly
                                        />
                                    </div>
                                    <div class="form-control">
                                        <label
                                            class="label py-2"
                                            for="longitude"
                                        >
                                            <span
                                                class="label-text font-bold text-base-content/80"
                                                >Longitude <span
                                                    class="text-error">*</span
                                                ></span
                                            >
                                        </label>
                                        <input
                                            id="longitude"
                                            type="number"
                                            step="any"
                                            class="input input-bordered w-full rounded-xl bg-base-200/50 min-h-12 cursor-not-allowed {validationErrors.longitude
                                                ? 'input-error'
                                                : ''}"
                                            bind:value={formData.longitude}
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
                                    use:mapAction={formData}
                                    class="relative z-0 h-72 w-full overflow-hidden rounded-2xl border border-base-200 shadow-inner sm:h-80 md:h-96"
                                ></div>
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

                        <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                            <div class="form-control">
                                <label class="label pt-0 pb-2" for="luas_tanah">
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Luas Tanah (m²) <span
                                            class="text-error">*</span
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
                                <label
                                    class="label pt-0 pb-2"
                                    for="luas_bangunan"
                                >
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Luas Bangunan (m²) <span
                                            class="text-error">*</span
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
                                <label class="label py-2" for="status_tanah">
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
                                    <option value="Sewa">Sewa</option>
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
                                <label
                                    class="label py-2"
                                    for="daya_tampung_max"
                                >
                                    <span
                                        class="label-text font-bold text-base-content/80"
                                        >Daya Tampung Max <span
                                            class="text-error">*</span
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
                                        (validationErrors.daya_tampung_max =
                                            null)}
                                />
                                {#if validationErrors.daya_tampung_max}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.daya_tampung_max}</span
                                    >
                                {/if}
                            </div>

                            <div class="form-control">
                                <label class="label py-2" for="nomor_aiw">
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
                                />
                                {#if validationErrors.nomor_aiw}
                                    <span
                                        class="mt-1.5 text-xs font-medium text-error"
                                        >{validationErrors.nomor_aiw}</span
                                    >
                                {/if}
                            </div>
                            <div class="form-control">
                                <label
                                    class="label py-2"
                                    for="nomor_sertifikat_wakaf"
                                >
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

                        {#if existingPhotos.length > 0}
                            <div class="mb-6">
                                <label
                                    for="existing-photos"
                                    class="label pt-0 pb-2"
                                >
                                    <span
                                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                                        >Foto Tersimpan</span
                                    >
                                </label>
                                <div
                                    class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4"
                                >
                                    {#each formData.foto || [] as photo (photo.id)}
                                        <div
                                            class="relative overflow-hidden rounded-xl border border-base-200 shadow-sm group"
                                        >
                                            <img
                                                src={photo.image_url}
                                                alt="Foto situs"
                                                class="h-32 w-full object-cover transition-transform duration-500 group-hover:scale-105"
                                            />
                                            <button
                                                type="button"
                                                class="btn btn-circle btn-xs btn-error absolute top-2 right-2 text-white shadow-lg opacity-90 hover:opacity-100"
                                                onclick={() =>
                                                    markPhotoForDeletion(
                                                        photo.id,
                                                    )}
                                                aria-label="Hapus foto"
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
                            </div>
                        {/if}

                        {#if newPhotos.length > 0}
                            <div class="mb-6">
                                <label for="new-photos" class="label pt-0 pb-2">
                                    <span
                                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                                        >Foto Baru (Akan Diunggah)</span
                                    >
                                </label>
                                <div
                                    class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4"
                                >
                                    {#each newPhotos as photo, index (photo.preview)}
                                        <div
                                            class="relative overflow-hidden rounded-xl border border-base-200 shadow-sm group"
                                        >
                                            <img
                                                src={photo.preview}
                                                alt="Preview"
                                                class="h-32 w-full object-cover opacity-70 transition-transform duration-500 group-hover:scale-105"
                                            />
                                            <button
                                                type="button"
                                                class="btn btn-circle btn-xs btn-error absolute top-2 right-2 text-white shadow-lg opacity-90 hover:opacity-100"
                                                onclick={() =>
                                                    removeNewPhoto(index)}
                                                aria-label="Batalkan unggah"
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
                                                class="absolute bottom-2 left-2 rounded-lg bg-black/70 backdrop-blur-sm px-2 py-1 text-xs font-bold text-white tracking-wider"
                                            >
                                                WebP
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/if}

                        <div class="form-control">
                            <label class="label pt-0 pb-2" for="images">
                                <span
                                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                                    >Unggah Foto Tambahan</span
                                >
                            </label>
                            <div class="flex flex-col gap-3">
                                <input
                                    id="images"
                                    type="file"
                                    accept="image/*"
                                    multiple
                                    class="file-input file-input-bordered w-full rounded-xl bg-base-200/30 transition-colors min-h-12"
                                    onchange={handlePhotoSelect}
                                    disabled={isConverting}
                                    bind:this={fileInputRef}
                                />
                                {#if isConverting}
                                    <div
                                        class="flex items-center gap-3 text-sm text-base-content/70 bg-base-200/50 rounded-xl p-4 font-medium border border-base-200 mt-2"
                                    >
                                        <span
                                            class="loading loading-sm loading-spinner text-primary"
                                        ></span>
                                        <span
                                            >Mengkonversi gambar ke format WebP
                                            untuk menghemat penyimpanan...</span
                                        >
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>

                <div
                    class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8 pt-4"
                >
                    <button
                        type="button"
                        class="btn btn-ghost rounded-xl px-8 min-h-12 font-bold text-base-content/60 hover:bg-base-200/50 hover:text-base-content w-full sm:w-auto"
                        onclick={goBack}
                        disabled={isSubmitting}
                    >
                        Batal
                    </button>
                    <button
                        type="submit"
                        class="btn btn-primary rounded-xl px-10 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-12 font-extrabold text-base-100 w-full sm:w-auto"
                        disabled={isSubmitting || isConverting}
                    >
                        {#if isSubmitting}
                            <span class="loading loading-spinner mr-2"></span> Menyimpan...
                        {:else}
                            Simpan Perubahan
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    {/if}
</div>

<style>
    :global(.leaflet-container) {
        z-index: 10;
    }
    :global(.leaflet-control-container) {
        z-index: 10;
    }
</style>
