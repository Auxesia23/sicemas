<script>
    const features = [
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>',
            title: "Manajemen Data Lengkap",
            description:
                "Kelola profil, kegiatan, dan aset masjid dengan sistem terintegrasi yang mudah digunakan.",
        },
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>',
            title: "Statistik Real-time",
            description:
                "Pantau perkembangan data situs keagamaan dengan dashboard analitik yang informatif.",
        },
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/></svg>',
            title: "Pemetaan Lokasi",
            description:
                "Lihat distribusi geografis masjid dan mushola di peta interaktif wilayah Ciemas.",
        },
    ];
</script>

<svelte:head>
    <title>Direktori Situs Keagamaan - KUA Ciemas | SICEMAS</title>
    <meta
        name="description"
        content="Sistem Informasi Catatan Elektronik Masjid (SICEMAS) Kecamatan Ciemas. Database lengkap masjid, mushola, pondok pesantren, dan majelis taklim di Sukabumi dengan data terupdate."
    />
    <meta
        name="keywords"
        content="masjid ciemas, pondok pesantren ciemas, majelis taklim ciemas, kua ciemas, situs keagamaan sukabumi, sicemas, mushola ciemas, data masjid sukabumi"
    />

    <meta
        property="og:description"
        content="Platform digital database masjid, mushola, dan pondok pesantren di wilayah Ciemas, Sukabumi."
    />
</svelte:head>

<a
    href="#main-content"
    class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 bg-primary text-white px-4 py-2 rounded z-50"
>
    Lompat ke konten utama
</a>

<main id="main-content">
    <section
        class="relative flex items-center pt-28 pb-20 md:pt-36 md:pb-28 overflow-hidden bg-linear-to-br from-base-100 via-base-100 to-primary/5"
    >
        <div class="absolute inset-0 opacity-30" aria-hidden="true">
            <div
                class="absolute inset-0"
                style="background-image: radial-gradient(circle at 1px 1px, hsl(var(--p)/0.15) 1px, transparent 0); background-size: 40px 40px;"
            ></div>
        </div>

        <div
            class="absolute inset-0 pointer-events-none overflow-hidden"
            aria-hidden="true"
        >
            <div
                class="animate-blob absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl"
            ></div>
            <div
                class="animate-blob animation-delay-2000 absolute bottom-1/4 right-1/4 w-96 h-96 bg-secondary/20 rounded-full blur-3xl"
            ></div>
            <div
                class="animate-blob animation-delay-4000 absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-125 h-125 bg-accent/10 rounded-full blur-3xl"
            ></div>
        </div>

        <div class="container mx-auto px-4 relative z-10">
            <div class="max-w-4xl mx-auto text-center">
                <div
                    class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-bold mb-6 border border-primary/20 backdrop-blur-sm"
                >
                    <span class="tracking-wide">KUA Kecamatan Ciemas</span>
                </div>

                <h1
                    class="text-4xl md:text-5xl lg:text-6xl font-bold text-base-content leading-tight mb-6"
                >
                    Sistem Informasi Catatan Elektronik Masjid Area Ciemas
                </h1>

                <p
                    class="text-lg md:text-xl text-base-content/70 max-w-2xl mx-auto mb-8 leading-relaxed"
                >
                    Platform digital terintegrasi untuk pengelolaan data masjid,
                    mushola, pondok pesantren, dan majelis taklim di wilayah
                    Kecamatan Ciemas, Kabupaten Sukabumi.
                </p>

                <div
                    class="flex flex-col sm:flex-row gap-4 justify-center items-center"
                >
                    <a
                        href="/dashboard"
                        class="btn btn-primary btn-lg px-8 gap-2 shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all"
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
                        Mulai Mengelola Data
                    </a>
                    <a
                        href="/situs"
                        class="btn btn-outline btn-lg px-8 gap-2 hover:bg-base-200 transition-all"
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
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                            />
                        </svg>
                        Jelajahi Direktori
                    </a>
                </div>
            </div>
        </div>
    </section>

    <section class="py-20 bg-base-200" aria-labelledby="features-heading">
        <div class="container mx-auto px-4">
            <div class="text-center max-w-2xl mx-auto mb-16">
                <h2
                    id="features-heading"
                    class="text-3xl md:text-4xl font-bold mb-4"
                >
                    Fitur Utama SICEMAS
                </h2>
                <p class="text-base-content/70">
                    Kelola data situs keagamaan dengan lebih efisien dan
                    terstruktur
                </p>
            </div>

            <div class="grid md:grid-cols-3 gap-8">
                {#each features as feature}
                    <article
                        class="bg-base-100 rounded-2xl p-6 shadow-sm hover:shadow-md transition-shadow border border-base-300"
                    >
                        <div
                            class="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mb-4 text-primary"
                        >
                            {@html feature.icon}
                        </div>
                        <h3 class="text-xl font-bold mb-2">{feature.title}</h3>
                        <p class="text-base-content/70 leading-relaxed">
                            {feature.description}
                        </p>
                    </article>
                {/each}
            </div>
        </div>
    </section>

    <section
        class="py-20 md:py-28 bg-base-100 relative overflow-hidden"
        aria-labelledby="about-heading"
    >
        <div
            class="absolute top-0 right-0 w-1/3 h-full bg-linear-to-l from-primary/5 to-transparent pointer-events-none"
            aria-hidden="true"
        ></div>
        <div
            class="absolute bottom-0 left-0 w-64 h-64 bg-secondary/10 rounded-full blur-3xl pointer-events-none"
            aria-hidden="true"
        ></div>

        <div class="container mx-auto px-4 relative z-10">
            <div
                class="flex flex-col lg:flex-row items-center gap-12 lg:gap-16"
            >
                <div class="lg:w-1/2 order-2 lg:order-1">
                    <div
                        class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 text-primary text-sm font-medium mb-6"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-4 w-4"
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
                        Tentang Kami
                    </div>

                    <h2
                        id="about-heading"
                        class="text-3xl md:text-4xl lg:text-5xl font-bold text-base-content mb-6 leading-tight"
                    >
                        KUA Kecamatan Ciemas
                    </h2>

                    <div
                        class="space-y-4 text-base-content/70 text-lg leading-relaxed mb-8"
                    >
                        <p>
                            Kantor Urusan Agama (KUA) Kecamatan Ciemas bertugas
                            melaksanakan tugas pemerintahan bidang urusan agama
                            di wilayah Kecamatan Ciemas, Kabupaten Sukabumi.
                        </p>
                        <p>
                            Dalam rangka mendukung pengelolaan situs keagamaan
                            yang ada di wilayah kerja kami, kami menghadirkan
                            sistem pendataan digital ini sebagai sarana pemetaan
                            dan monitoring yang efektif dan efisien.
                        </p>
                    </div>

                    <div class="grid grid-cols-2 gap-4 mb-8">
                        <div
                            class="flex items-start gap-3 p-4 rounded-xl bg-base-200/50 border border-base-200"
                        >
                            <div
                                class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0 text-primary"
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
                                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </div>
                            <div>
                                <h4 class="font-semibold text-base-content">
                                    Terintegrasi
                                </h4>
                                <p class="text-sm text-base-content/60">
                                    Satu platform untuk semua data
                                </p>
                            </div>
                        </div>
                        <div
                            class="flex items-start gap-3 p-4 rounded-xl bg-base-200/50 border border-base-200"
                        >
                            <div
                                class="w-10 h-10 rounded-lg bg-secondary/10 flex items-center justify-center shrink-0 text-secondary"
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
                                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </div>
                            <div>
                                <h4 class="font-semibold text-base-content">
                                    Real-time
                                </h4>
                                <p class="text-sm text-base-content/60">
                                    Update data instan
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="lg:w-1/2 order-1 lg:order-2 w-full">
                    <div class="relative">
                        <div
                            class="absolute -top-4 -right-4 w-24 h-24 bg-primary/20 rounded-full blur-2xl"
                            aria-hidden="true"
                        ></div>
                        <div
                            class="absolute -bottom-4 -left-4 w-32 h-32 bg-secondary/20 rounded-full blur-2xl"
                            aria-hidden="true"
                        ></div>

                        <div
                            class="relative rounded-2xl overflow-hidden shadow-2xl bg-linear-to-br from-primary to-secondary p-1"
                        >
                            <div class="bg-base-100 rounded-xl overflow-hidden">
                                <iframe
                                    title="Lokasi KUA Kecamatan Ciemas di Google Maps"
                                    src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3958.02010347061!2d106.51973537553341!3d-7.238545592767754!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e68200002b8911b%3A0xae335d95d6fefa1d!2sKUA%20KEC.%20CIEMAS!5e0!3m2!1sen!2sid!4v1773645785108!5m2!1sen!2sid"
                                    width="100%"
                                    height="450"
                                    style="border:0; filter: grayscale(0.1) contrast(1.05);"
                                    loading="lazy"
                                    referrerpolicy="no-referrer-when-downgrade"
                                    class="w-full"
                                ></iframe>

                                <div
                                    class="p-4 bg-base-100 border-t border-base-200"
                                >
                                    <div class="flex items-start gap-3">
                                        <div
                                            class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0 text-primary"
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
                                        <div>
                                            <h4
                                                class="font-semibold text-base-content"
                                            >
                                                Alamat Kantor
                                            </h4>
                                            <p
                                                class="text-sm text-base-content/70"
                                            >
                                                Jl. Raya Ciemas No. 123,
                                                Kecamatan Ciemas, Kabupaten
                                                Sukabumi, Jawa Barat 43177
                                            </p>
                                            <a
                                                href="https://maps.google.com/?q=KUA+KEC.+CIEMAS"
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                class="inline-flex items-center gap-1 text-sm text-primary hover:underline mt-1"
                                            >
                                                Buka di Google Maps
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
                                                        d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                                                    />
                                                </svg>
                                            </a>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
    <section
        class="py-20 bg-primary text-primary-content relative overflow-hidden"
    >
        <div class="absolute inset-0 opacity-10" aria-hidden="true">
            <div
                class="absolute inset-0"
                style="background-image: radial-gradient(circle at 2px 2px, white 1px, transparent 0); background-size: 32px 32px;"
            ></div>
        </div>

        <div class="container mx-auto px-4 relative z-10 text-center">
            <h2 class="text-3xl md:text-4xl font-bold mb-4">
                Siap Mengelola Data Situs Keagamaan?
            </h2>
            <p class="text-primary-content/80 max-w-2xl mx-auto mb-8 text-lg">
                Bergabung dengan SICEMAS untuk mendapatkan akses penuh ke
                dashboard pengelolaan data dan laporan statistik real-time.
            </p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
                <a href="/dashboard" class="btn btn-secondary btn-lg px-8"
                    >Dashboard</a
                >
            </div>
        </div>
    </section>
</main>

<style>
    /* Custom animations */
    @keyframes blob {
        0%,
        100% {
            transform: translate(0, 0) scale(1);
        }
        33% {
            transform: translate(30px, -50px) scale(1.1);
        }
        66% {
            transform: translate(-20px, 20px) scale(0.9);
        }
    }

    .animate-blob {
        animation: blob 7s infinite;
    }

    .animation-delay-2000 {
        animation-delay: 2s;
    }

    .animation-delay-4000 {
        animation-delay: 4s;
    }

    :focus-visible {
        outline: 2px solid hsl(var(--p));
        outline-offset: 2px;
    }
</style>
