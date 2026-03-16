<script>
	let { data } = $props();
	let stats = $derived(
		data.stats || { total_situs: 0, total_desa: 0, total_kategori: 0, total_kapasitas: 0 }
	);
	let user = $derived(data.user);
	function formatAngkaBesar(num) {
		if (!num) return '0';

		if (num >= 1000000) {
			return (num / 1000000).toFixed(1).replace('.', ',') + ' Juta';
		}
		if (num >= 10000) {
			return (num / 1000).toFixed(1).replace('.', ',') + ' Ribu';
		}

		return num.toLocaleString('id-ID');
	}
</script>

<div class="min-h-screen bg-linear-to-b from-base-100 to-base-200">
	<!-- Navigation -->
	<nav class="navbar bg-base-100 px-4 sm:px-6 lg:px-8">
		<div class="flex-1">
			<a href="/" class="btn text-xl font-bold text-primary btn-ghost">SICEMAS</a>
		</div>
		<div class="flex-none">
			<a href="/login" class="btn btn-primary"
				>{#if user}Dashboard{:else}Masuk{/if}</a
			>
		</div>
	</nav>

	<!-- Hero Section -->
	<section class="min-h-screen py-16 md:py-24">
		<div class="container mx-auto px-4 text-center">
			<div class="mx-auto max-w-3xl">
				<span class="mb-4 badge badge-primary">KUA Kecamatan Ciemas</span>
				<h1 class="mb-6 text-4xl leading-tight font-bold text-gray-800 md:text-6xl">
					Sistem Informasi Catatan Elektronik <span
						class="bg-linear-to-r from-primary to-secondary bg-clip-text text-transparent"
						>Masjid Area Ciemas</span
					>
				</h1>
				<p class="mx-auto mb-10 max-w-2xl text-lg text-gray-600 md:text-xl">
					Platform digital untuk mengelola dan memantau situs-situs keagamaan di wilayah Kecamatan
					Ciemas, Kabupaten Sukabumi
				</p>
				<div class="flex flex-col justify-center gap-4 sm:flex-row">
					<a href="/dashboard" class="btn px-8 btn-lg btn-primary">Mulai Mengelola Data</a>
					<a href="/situs" class="btn px-8 btn-outline btn-lg">Jelajahi Situs</a>
				</div>
			</div>

			<!-- Decorative Elements -->
			<div class="relative mt-16">
				<div class="absolute inset-0 flex items-center justify-center">
					<div
						class="animate-blob filter h-64 w-64 rounded-full bg-primary opacity-20 mix-blend-multiply blur-xl"
					></div>
					<div
						class="animate-blob animation-delay-2000 filter ml-64 h-64 w-64 rounded-full bg-secondary opacity-20 mix-blend-multiply blur-xl"
					></div>
					<div
						class="animate-blob animation-delay-4000 filter mr-64 h-64 w-64 rounded-full bg-accent opacity-20 mix-blend-multiply blur-xl"
					></div>
				</div>
				<div class="relative">
					<div
						class="mx-auto max-w-4xl rounded-2xl border border-base-200 bg-base-100 p-6 shadow-xl"
					>
						<div class="grid grid-cols-2 gap-6 text-center md:grid-cols-4">
							<div class="p-4">
								<div class="text-3xl font-bold text-primary">{stats.total_situs}</div>
								<div class="text-sm text-gray-500">Situs Terdata</div>
							</div>

							<div class="p-4">
								<div class="text-3xl font-bold text-secondary">{stats.total_desa}</div>
								<div class="text-sm text-gray-500">Desa Terjangkau</div>
							</div>

							<div class="p-4">
								<div class="text-3xl font-bold text-accent">{stats.total_kategori}</div>
								<div class="text-sm text-gray-500">Kategori Situs</div>
							</div>

							<div class="p-4">
								<div class="text-3xl font-bold text-primary">
									{formatAngkaBesar(stats.total_kapasitas)}
								</div>
								<div class="text-sm text-gray-500">Kapasitas Jamaah</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="bg-base-100 py-16">
		<div class="container mx-auto px-4">
			<div class="mb-16 text-center">
				<h2 class="mb-4 text-3xl font-bold text-gray-800 md:text-4xl">Fitur Unggulan</h2>
				<p class="mx-auto max-w-2xl text-gray-600">
					Sistem kami dirancang untuk memudahkan pengelolaan data situs keagamaan secara
					komprehensif
				</p>
			</div>

			<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
				<!-- Feature 1: Pencarian Pintar -->
				<div
					class="card border border-base-200 bg-base-100 shadow-lg transition-shadow duration-300 hover:shadow-xl"
				>
					<div class="card-body p-6">
						<div
							class="bg-opacity-10 mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-primary"
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
						</div>
						<h3 class="card-title text-lg font-bold text-gray-800">Pencarian Pintar</h3>
						<p class="text-gray-600">
							Cari masjid, majelis taklim, atau pesantren berdasarkan nama, jenis, atau nama desa
							dengan cepat.
						</p>
					</div>
				</div>

				<!-- Feature 2: Detail Fasilitas & Galeri -->
				<div
					class="card border border-base-200 bg-base-100 shadow-lg transition-shadow duration-300 hover:shadow-xl"
				>
					<div class="card-body p-6">
						<div
							class="bg-opacity-10 mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-secondary"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-secondary"
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
						<h3 class="card-title text-lg font-bold text-gray-800">Detail Fasilitas & Galeri</h3>
						<p class="text-gray-600">
							Lihat kelengkapan fasilitas, kapasitas jamaah, luas tanah, hingga foto-foto kondisi
							situs terkini.
						</p>
					</div>
				</div>

				<!-- Feature 3: Integrasi Navigasi -->
				<div
					class="card border border-base-200 bg-base-100 shadow-lg transition-shadow duration-300 hover:shadow-xl"
				>
					<div class="card-body p-6">
						<div
							class="bg-opacity-10 mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-accent"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-accent"
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
						<h3 class="card-title text-lg font-bold text-gray-800">Integrasi Navigasi</h3>
						<p class="text-gray-600">
							Dapatkan titik koordinat akurat dan panduan rute langsung menggunakan integrasi Google
							Maps.
						</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Stats Section -->
	<section class="bg-linear-to-r from-primary to-secondary py-16 text-primary-content">
		<div class="container mx-auto px-4">
			<div class="grid grid-cols-2 gap-8 text-center md:grid-cols-4">
				<div>
					<div class="text-4xl font-bold md:text-5xl">{stats.total_situs}</div>
					<div class="text-sm opacity-80">Situs Terdata</div>
				</div>
				<div>
					<div class="text-4xl font-bold md:text-5xl">{stats.total_desa}</div>
					<div class="text-sm opacity-80">Desa Terjangkau</div>
				</div>
				<div>
					<div class="text-4xl font-bold md:text-5xl">{stats.total_kategori}</div>
					<div class="text-sm opacity-80">Kategori Situs</div>
				</div>
				<div>
					<div class="text-4xl font-bold md:text-5xl">
						{formatAngkaBesar(stats.total_kapasitas)}
					</div>
					<div class="text-sm opacity-80">Kapasitas Jemaah</div>
				</div>
			</div>
		</div>
	</section>

	<!-- About Section -->
	<section class="bg-base-100 py-16">
		<div class="container mx-auto px-4">
			<div class="flex flex-col items-center gap-12 lg:flex-row">
				<div class="lg:w-1/2">
					<h2 class="mb-4 text-3xl font-bold text-gray-800">Tentang KUA Kecamatan Ciemas</h2>
					<p class="mb-4 text-gray-600">
						Kantor Urusan Agama (KUA) Kecamatan Ciemas bertugas melaksanakan tugas pemerintahan
						bidang urusan agama di wilayah Kecamatan Ciemas, Kabupaten Sukabumi.
					</p>
					<p class="mb-6 text-gray-600">
						Dalam rangka mendukung pengelolaan situs keagamaan yang ada di wilayah kerja kami, kami
						menghadirkan sistem pendataan digital ini sebagai sarana pemetaan dan monitoring yang
						efektif dan efisien.
					</p>
					<div class="flex flex-wrap gap-2">
						<span class="badge badge-outline">Digitalisasi</span>
						<span class="badge badge-outline">Transparansi</span>
						<span class="badge badge-outline">Efisiensi</span>
						<span class="badge badge-outline">Akuntabilitas</span>
					</div>
				</div>
				<div class="flex justify-center lg:w-1/2">
					<div
						class="w-full max-w-md rounded-2xl bg-linear-to-br from-primary to-secondary p-1 shadow-xl"
					>
						<iframe
							title="KUA Ciemas"
							src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3958.02010347061!2d106.51973537553341!3d-7.238545592767754!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e68200002b8911b%3A0xae335d95d6fefa1d!2sKUA%20KEC.%20CIEMAS!5e0!3m2!1sen!2sid!4v1773645785108!5m2!1sen!2sid"
							width="100%"
							height="400"
							style="border:0;"
							allowfullscreen=""
							loading="lazy"
							referrerpolicy="no-referrer-when-downgrade"
							class="min-h-75 rounded-xl"
						></iframe>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="py-16">
		<div class="container mx-auto px-4 text-center">
			<div class="mx-auto max-w-3xl">
				<h2 class="mb-6 text-3xl font-bold text-gray-800 md:text-4xl">
					Siap Kelola Data Situs Keagamaan Anda?
				</h2>
				<p class="mx-auto mb-10 max-w-2xl text-gray-600">
					Bergabunglah dengan sistem pendataan digital kami untuk meningkatkan efisiensi pengelolaan
					situs-situs keagamaan di wilayah Kecamatan Ciemas
				</p>
				<a href="/dashboard" class="btn px-10 btn-lg btn-primary">Akses Dashboard Sekarang</a>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="footer-center footer rounded bg-base-200 p-10 text-base-content">
		<div>
			<div class="grid grid-flow-col gap-4">
				<a
					href="https://www.tiktok.com/@kuaciemas?is_from_webapp=1&sender_device=pc"
					aria-label="TikTok"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						class="fill-current"
					>
						<path
							d="M12.53.02C13.84 0 15.14.01 16.44 0c.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z"
						></path>
					</svg>
				</a>
				<a href="https://www.instagram.com/kuaciemas_official/" aria-label="Instagram">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						class="fill-current"
					>
						<path
							d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z"
						></path>
					</svg>
				</a>
			</div>
		</div>
	</footer>
</div>

<style>
	@keyframes blob {
		0% {
			transform: translate(0px, 0px) scale(1);
		}
		33% {
			transform: translate(30px, -50px) scale(1.1);
		}
		66% {
			transform: translate(-20px, 20px) scale(0.9);
		}
		100% {
			transform: translate(0px, 0px) scale(1);
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
</style>
