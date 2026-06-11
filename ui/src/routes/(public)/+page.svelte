<script>
    const features = [
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>',
            title: "Pendataan Keagamaan",
            description:
                "Basis data profil, kegiatan, dan aset masjid serta mushola yang dikelola secara terpusat.",
        },
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>',
            title: "Laporan & Statistik",
            description:
                "Rekapitulasi data situs keagamaan yang diperbarui secara berkala untuk kebutuhan administrasi.",
        },
        {
            icon: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/></svg>',
            title: "Peta Lokasi",
            description:
                "Pemetaan geografis untuk mempermudah pencarian lokasi masjid dan pesantren di wilayah Ciemas.",
        },
    ];
</script>

<svelte:head>
    <title>Direktori Situs Keagamaan - KUA Ciemas | SICEMAS</title>
    <meta
        name="description"
        content="Sistem Informasi Catatan Elektronik Masjid (SICEMAS) Kecamatan Ciemas. Database masjid, mushola, pondok pesantren, dan majelis taklim di Sukabumi."
    />
    <meta
        name="keywords"
        content="masjid ciemas, pondok pesantren ciemas, majelis taklim ciemas, kua ciemas, situs keagamaan sukabumi, sicemas, mushola ciemas"
    />
    <meta
        property="og:description"
        content="Pusat data masjid, mushola, dan pondok pesantren di wilayah Ciemas, Sukabumi."
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
        class="min-h-screen flex items-center justify-center py-20 bg-base-100 border-b border-base-200"
    >
        <div class="container mx-auto px-4">
            <div class="max-w-3xl mx-auto text-center">
                <span
                    class="text-primary font-bold text-sm uppercase tracking-wider mb-6 block"
                >
                    Kantor Urusan Agama Kecamatan Ciemas
                </span>

                <h1
                    class="text-4xl md:text-5xl lg:text-6xl font-extrabold text-base-content leading-tight mb-8"
                >
                    Sistem Informasi Catatan Elektronik Masjid (SICEMAS)
                </h1>

                <p
                    class="text-lg md:text-xl text-base-content/70 mb-10 leading-relaxed"
                >
                    Layanan pendataan resmi untuk masjid, mushola, pondok
                    pesantren, dan majelis taklim di wilayah Kecamatan Ciemas,
                    Kabupaten Sukabumi.
                </p>

                <div class="flex justify-center">
                    <a href="/dashboard" class="btn btn-primary btn-lg px-10">
                        Masuk ke Dashboard
                    </a>
                </div>
            </div>
        </div>
    </section>

    <section
        class="min-h-screen flex items-center py-20 bg-base-200"
        aria-labelledby="features-heading"
    >
        <div class="container mx-auto px-4">
            <div class="text-center max-w-2xl mx-auto mb-16">
                <h2
                    id="features-heading"
                    class="text-3xl md:text-4xl font-bold mb-4 text-base-content"
                >
                    Layanan Pendataan
                </h2>
                <p class="text-lg text-base-content/70">
                    Fasilitas sistem informasi untuk pengelolaan data keagamaan
                    wilayah Ciemas.
                </p>
            </div>

            <div class="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
                {#each features as feature}
                    <article
                        class="bg-base-100 rounded-xl p-8 border border-base-300 shadow-sm hover:shadow-md transition-shadow"
                    >
                        <div
                            class="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-6 text-primary"
                        >
                            {@html feature.icon}
                        </div>
                        <h3 class="text-xl font-bold mb-3 text-base-content">
                            {feature.title}
                        </h3>
                        <p class="text-base-content/70 leading-relaxed">
                            {feature.description}
                        </p>
                    </article>
                {/each}
            </div>
        </div>
    </section>

    <section
        class="min-h-screen flex items-center py-20 bg-base-100 border-t border-base-200"
        aria-labelledby="about-heading"
    >
        <div class="container mx-auto px-4 max-w-6xl">
            <div
                class="flex flex-col lg:flex-row gap-12 lg:gap-16 items-center"
            >
                <div class="lg:w-1/2">
                    <h2
                        id="about-heading"
                        class="text-3xl md:text-4xl font-bold text-base-content mb-6"
                    >
                        Tentang Program SICEMAS
                    </h2>

                    <div
                        class="space-y-5 text-base-content/70 leading-relaxed mb-8 text-lg"
                    >
                        <p>
                            Kantor Urusan Agama (KUA) Kecamatan Ciemas
                            menginisiasi Sistem Informasi Catatan Elektronik
                            Masjid (SICEMAS) untuk mewujudkan tertib
                            administrasi pada lembaga keagamaan di wilayah kerja
                            kami.
                        </p>
                        <p>
                            Sistem ini dirancang untuk memudahkan pemantauan,
                            pembinaan, serta penyaluran informasi program dari
                            pemerintah pusat maupun daerah secara akurat dan
                            tepat sasaran.
                        </p>
                    </div>

                    <div
                        class="bg-base-200 border border-base-300 rounded-xl p-6"
                    >
                        <h4 class="font-bold text-base-content mb-2">
                            Kantor Layanan
                        </h4>
                        <p class="text-base-content/70 mb-4">
                            Jl. Raya Ciemas No. 123, Kecamatan Ciemas, Kabupaten
                            Sukabumi, Jawa Barat 43177
                        </p>
                        <a
                            href="https://maps.google.com/?q=KUA+KEC.+CIEMAS"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="text-primary hover:text-primary-focus hover:underline font-semibold inline-flex items-center gap-1"
                        >
                            Buka di Google Maps →
                        </a>
                    </div>
                </div>

                <div class="lg:w-1/2 w-full">
                    <div
                        class="rounded-xl border border-base-300 overflow-hidden bg-base-200 shadow-sm"
                    >
                        <iframe
                            title="Lokasi KUA Kecamatan Ciemas di Google Maps"
                            src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3958.02010347061!2d106.51973537553341!3d-7.238545592767754!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e68200002b8911b%3A0xae335d95d6fefa1d!2sKUA%20KEC.%20CIEMAS!5e0!3m2!1sen!2sid!4v1773645785108!5m2!1sen!2sid"
                            width="100%"
                            height="450"
                            style="border:0;"
                            loading="lazy"
                            referrerpolicy="no-referrer-when-downgrade"
                            class="w-full block"
                        ></iframe>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <section
        class="py-24 bg-primary text-primary-content flex items-center justify-center"
    >
        <div class="container mx-auto px-4 text-center">
            <h2 class="text-3xl md:text-4xl font-bold mb-6">Akses Pengelola</h2>
            <p class="text-primary-content/90 max-w-xl mx-auto mb-10 text-lg">
                Silakan masuk menggunakan kredensial yang telah didaftarkan
                untuk mengelola data lembaga keagamaan Anda secara mandiri.
            </p>
            <a
                href="/dashboard"
                class="btn btn-lg bg-base-100 text-primary hover:bg-base-200 border-none px-12 shadow-lg"
            >
                Login Sistem
            </a>
        </div>
    </section>
</main>

<style>
    :focus-visible {
        outline: 2px solid hsl(var(--p));
        outline-offset: 2px;
    }
</style>
