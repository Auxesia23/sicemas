<script lang="ts">
    import type { PageData } from "./$types";

    let { data }: { data: PageData } = $props();

    let user = $derived(data.profile);

    function formatDate(isoString: string | undefined | null) {
        if (!isoString) return "-";
        const date = new Date(isoString);
        return new Intl.DateTimeFormat("id-ID", {
            dateStyle: "long",
        }).format(date);
    }

    function getInitials(name = "") {
        if (!name) return "??";
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    }
</script>

<svelte:head>
    <title>Profil Pengguna | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
</svelte:head>

<div class="w-full">
    <div
        class="mb-8 flex flex-col gap-6 md:flex-row md:items-end md:justify-between"
    >
        <div>
            <div class="flex items-center gap-3 mb-2">
                <div
                    class="w-1.5 h-8 rounded-full bg-linear-to-b from-primary to-secondary"
                ></div>
                <h1
                    class="text-3xl font-extrabold tracking-tight text-base-content"
                >
                    Profil Pengguna
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Informasi detail akun dan identitas Anda di sistem SICEMAS
            </p>
        </div>
    </div>

    <div
        class="card bg-base-100 shadow-xl shadow-base-200/40 border border-base-200 rounded-2xl overflow-hidden"
    >
        <div
            class="h-32 bg-linear-to-r from-primary/80 via-primary to-secondary/80 w-full relative"
        >
            <div
                class="absolute inset-0 opacity-20"
                style="background-image: radial-gradient(circle at 2px 2px, white 1px, transparent 0); background-size: 24px 24px;"
            ></div>
        </div>

        <div class="card-body pt-0 px-6 sm:px-10 pb-10 relative">
            <div
                class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-10"
            >
                <div class="avatar -mt-16 z-10">
                    <div
                        class="h-32 w-32 rounded-2xl ring-8 ring-base-100 bg-base-200 shadow-xl flex items-center justify-center overflow-hidden bg-linear-to-br from-primary/10 to-secondary/10"
                    >
                        <div
                            class="w-full h-full flex items-center justify-center text-3xl font-black text-primary"
                        >
                            {getInitials(user.nama_lengkap)}
                        </div>
                    </div>
                </div>

                <div class="flex pb-2">
                    <div
                        class="badge bg-success/10 text-success border-success/20 gap-2 p-4 font-bold rounded-xl shadow-sm"
                    >
                        <div
                            class="w-2 h-2 rounded-full bg-success animate-pulse"
                        ></div>
                        Akun Aktif
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-primary/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary shrink-0 transition-transform group-hover:scale-110"
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
                                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Nama Lengkap
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {user.nama_lengkap}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-primary/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary shrink-0 transition-transform group-hover:scale-110"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/xl"
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            NIP / NIK
                        </p>
                        <p
                            class="text-base font-mono font-bold text-base-content truncate"
                        >
                            {user.nip}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-secondary/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-secondary/10 flex items-center justify-center text-secondary shrink-0 transition-transform group-hover:scale-110"
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
                                d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Jabatan
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {user.jabatan}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-secondary/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-secondary/10 flex items-center justify-center text-secondary shrink-0 transition-transform group-hover:scale-110"
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
                                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Alamat Email
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {user.email}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-accent/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center text-accent shrink-0 transition-transform group-hover:scale-110"
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
                                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Unit Kerja
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {user.unit_kerja}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-accent/30 group"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center text-accent shrink-0 transition-transform group-hover:scale-110"
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
                                d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Nomor Telepon
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {user.nomor_telepon}
                        </p>
                    </div>
                </div>

                <div
                    class="flex items-center gap-4 p-5 rounded-2xl bg-base-200/30 border border-base-200 transition-all hover:border-info/30 group md:col-span-2"
                >
                    <div
                        class="w-12 h-12 rounded-xl bg-info/10 flex items-center justify-center text-info shrink-0 transition-transform group-hover:scale-110"
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
                                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                            />
                        </svg>
                    </div>
                    <div class="min-w-0 flex-1">
                        <p
                            class="text-[10px] font-bold text-base-content/40 uppercase tracking-widest mb-0.5"
                        >
                            Bergabung Sejak
                        </p>
                        <p
                            class="text-base font-bold text-base-content truncate"
                        >
                            {formatDate(user.created_at)}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
