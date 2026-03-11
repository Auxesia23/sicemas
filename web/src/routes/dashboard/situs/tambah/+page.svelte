<script>
	import { onMount } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import apiService from '$lib/api';
	import { goto } from '$app/navigation';

	// State for form data
	let formData = $state({
		jenis_situs_id: '',
		nama: '',
		jenis_tipologi: '',
		nomor_telepon: '',
		email: '',
		website: '',
		nomor_badan_hukum: '',
		tahun_berdiri: '',
		provinsi: 'Jawa Barat',
		kabupaten_kota: 'Kab. Sukabumi',
		kecamatan: 'Ciemas',
		desa: '',
		alamat_lengkap: '',
		latitude: null,
		longitude: null,
		luas_tanah: null,
		luas_bangunan: null,
		status_tanah: '',
		nomor_aiw: '',
		nomor_sertifikat_wakaf: '',
		daya_tampung_max: null,
		detail: {}
	});

	// State for site types dropdown
	let jenisSitusList = $state([]);
	let selectedJenisSitusName = $state('');

	// State for map
	let mapContainer;
	let map;
	let marker;

	// State for dynamic list inputs (Masjid detail fields)
	let newJenisBuku = $state('');
	let newNamaImam = $state('');
	let newNamaMuazdin = $state('');
	let newFasilitasUmum = $state('');
	let newFasilitasRamanAnak = $state('');
	let newFasilitasDisabilitas = $state('');
	let newKegiatan = $state('');

	// State for toast notification
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success'); // 'success' or 'error'

	// Derived: check if selected site is Masjid
	let isMasjid = $derived(selectedJenisSitusName.toLowerCase().includes('masjid'));

	// Initialize detail object for Masjid
	function initMasjidDetail() {
		formData.detail = {
			perpustakaan: {
				luas_m2: null,
				jumlah_pengurus: null,
				kondisi: '',
				jenis_buku: []
			},
			kalibrasi_arah_kiblat: {
				azimut: '',
				tanggal_kalibrasi: ''
			},
			sdm_masjid: {
				jumlah_pengurus: null,
				jumlah_imam: null,
				jumlah_khotib: null,
				jumlah_muadzin: null,
				jumlah_remaja: null,
				jumlah_jemaah: null
			},
			pengurus_dkm: {
				ketua: '',
				sekretaris: '',
				bendahara: '',
				nama_imam: [],
				nama_muazdin: []
			},
			fasilitas_umum: [],
			fasilitas_raman_anak: [],
			fasilitas_disabilitas: [],
			kegiatan: []
		};
	}

	// Initialize map after component mounts
	onMount(async () => {
		try {
			const response = await apiService.get('/jenis-situs');
			const data = await response.json();
			jenisSitusList = data.data || data;
		} catch (error) {
			console.error('Failed to fetch jenis situs:', error);
		}
		if (mapContainer && !map) {
			const L = await import('leaflet');
			map = L.map(mapContainer).setView([-6.2088, 106.8456], 13);

			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				attribution: '© OpenStreetMap contributors'
			}).addTo(map);

			marker = L.marker([-6.2088, 106.8456], { draggable: true }).addTo(map);

			// Update form data when marker is dragged
			marker.on('dragend', (e) => {
				const latlng = e.target.getLatLng();
				formData.latitude = latlng.lat;
				formData.longitude = latlng.lng;
			});

			// Update marker when clicking on map
			map.on('click', (e) => {
				formData.latitude = e.latlng.lat;
				formData.longitude = e.latlng.lng;
				marker.setLatLng(e.latlng);
			});
		}

		return () => {
			if (map) {
				map.remove();
				map = null;
			}
		};
	});

	// Effect: Initialize Masjid detail when Masjid is selected
	$effect(() => {
		if (isMasjid && !formData.detail.perpustakaan) {
			initMasjidDetail();
		} else if (!isMasjid && Object.keys(formData.detail).length > 0) {
			formData.detail = {};
		}
	});

	// Effect: Update map marker when coordinates change
	$effect(() => {
		if (map && marker && formData.latitude !== null && formData.longitude !== null) {
			const newLatLng = new L.LatLng(formData.latitude, formData.longitude);
			marker.setLatLng(newLatLng);
			map.setView(newLatLng, 15);
		}
	});

	// Handle jenis situs change
	function handleJenisSitusChange(event) {
		const selectedId = event.target.value;
		formData.jenis_situs_id = selectedId;

		const selected = jenisSitusList.find((item) => item.id === selectedId);
		if (selected) {
			selectedJenisSitusName = selected.nama || selected.name || '';
			formData.jenis_tipologi = selectedJenisSitusName;
		}
	}

	// Get current location
	async function getCurrentLocation() {
		if (!navigator.geolocation) {
			alert('Geolocation is not supported by your browser');
			return;
		}

		try {
			const position = await new Promise((resolve, reject) => {
				navigator.geolocation.getCurrentPosition(resolve, reject, {
					enableHighAccuracy: true,
					timeout: 15000,
					maximumAge: 0
				});
			});

			formData.latitude = position.coords.latitude;
			formData.longitude = position.coords.longitude;
		} catch (error) {
			console.error('Geolocation error:', error);
			let errorMessage = 'Unable to get your location.';
			if (error.code === 1) {
				errorMessage = 'Location permission denied. Please allow location access and try again.';
			} else if (error.code === 2) {
				errorMessage =
					'Location unavailable. Please enter coordinates manually or click on the map.';
			} else if (error.code === 3) {
				errorMessage = 'Location request timed out. Please try again.';
			}
			alert(errorMessage);
		}
	}

	// Dynamic list helpers
	function addToArray(array, value, clearInput = true) {
		if (value.trim()) {
			array.push(value.trim());
			if (clearInput) return '';
		}
		return value;
	}

	function removeFromArray(array, index) {
		array.splice(index, 1);
	}

	function addJenisBuku() {
		newJenisBuku = addToArray(formData.detail.perpustakaan?.jenis_buku, newJenisBuku, true);
	}

	function addNamaImam() {
		newNamaImam = addToArray(formData.detail.pengurus_dkm?.nama_imam, newNamaImam, true);
	}

	function addNamaMuazdin() {
		newNamaMuazdin = addToArray(formData.detail.pengurus_dkm?.nama_muazdin, newNamaMuazdin, true);
	}

	function addFasilitasUmum() {
		newFasilitasUmum = addToArray(formData.detail.fasilitas_umum, newFasilitasUmum, true);
	}

	function addFasilitasRamanAnak() {
		newFasilitasRamanAnak = addToArray(
			formData.detail.fasilitas_raman_anak,
			newFasilitasRamanAnak,
			true
		);
	}

	function addFasilitasDisabilitas() {
		newFasilitasDisabilitas = addToArray(
			formData.detail.fasilitas_disabilitas,
			newFasilitasDisabilitas,
			true
		);
	}

	function addKegiatan() {
		newKegiatan = addToArray(formData.detail.kegiatan, newKegiatan, true);
	}

	// Handle enter key for list inputs
	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}

	// Submit handler
	async function handleSubmit() {
		try {
			const response = await apiService.post('/situs', formData);
			if (response.status === 201) {
				// Show success toast notification
				toastMessage = 'Data situs berhasil disimpan';
				toastType = 'success';
				showToast = true;
				setTimeout(() => {
					showToast = false;
					goto('/dashboard/situs');
				}, 3000);
			} else {
				toastMessage = 'Gagal menyimpan data situs';
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 5000);
			}
		} catch (error) {
			console.error('Failed to submit form:', error);
			toastMessage = 'Terjadi kesalahan. Silakan coba lagi.';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 5000);
		}
	}

	// Helper functions for detail field bindings
	function updatePerpustakaan(field, value) {
		if (!formData.detail.perpustakaan) initMasjidDetail();
		formData.detail.perpustakaan[field] = value;
	}

	function updateKalibrasi(field, value) {
		if (!formData.detail.kalibrasi_arah_kiblat) initMasjidDetail();
		formData.detail.kalibrasi_arah_kiblat[field] = value;
	}

	function updateSdm(field, value) {
		if (!formData.detail.sdm_masjid) initMasjidDetail();
		formData.detail.sdm_masjid[field] = value;
	}

	function updatePengurus(field, value) {
		if (!formData.detail.pengurus_dkm) initMasjidDetail();
		formData.detail.pengurus_dkm[field] = value;
	}
</script>

<svelte:head>
	<title>Tambah Situs Keagamaan</title>
</svelte:head>

<div class="min-h-screen bg-base-200 p-3 sm:p-4">
	<div class="mx-auto max-w-3xl">
		<!-- Header -->
		<div class="mb-4 sm:mb-6">
			<h1 class="text-2xl font-bold text-primary sm:text-3xl">Pendataan Situs Keagamaan</h1>
			<p class="mt-1 text-sm text-base-content/70 sm:mt-2 sm:text-base">
				Formulir pendaftaran situs keagamaan baru
			</p>
		</div>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-3 sm:space-y-4"
		>
			<!-- Jenis Situs Selection -->
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Jenis Situs</h2>
					<div class="form-control w-full">
						<label class="label" for="jenis_situs_id">
							<span class="label-text font-medium">Pilih Jenis Situs</span>
						</label>
						<select
							id="jenis_situs_id"
							class="select-bordered select min-h-11 w-full"
							value={formData.jenis_situs_id}
							oninput={handleJenisSitusChange}
						>
							<option required value="" disabled>Pilih jenis situs...</option>
							{#each jenisSitusList as jenis (jenis.id)}
								<option value={jenis.id}>{jenis.nama || jenis.name}</option>
							{/each}
						</select>
					</div>
				</div>
			</div>

			<!-- Informasi Umum -->
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Informasi Umum</h2>
					<div class="grid grid-cols-1 gap-3 sm:gap-4">
						<div class="form-control">
							<label class="label" for="nama">
								<span class="label-text font-medium">Nama Situs</span>
							</label>
							<input
								required
								id="nama"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="Masukkan nama situs"
								bind:value={formData.nama}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="jenis_tipologi">
								<span class="label-text font-medium">Jenis Tipologi</span>
							</label>
							<input
								required
								id="jenis_tipologi"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="Jenis tipologi"
								bind:value={formData.jenis_tipologi}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="nomor_telepon">
								<span class="label-text font-medium">Nomor Telepon</span>
							</label>
							<input
								id="nomor_telepon"
								type="tel"
								class="input-bordered input min-h-11 w-full"
								placeholder="081234567890"
								bind:value={formData.nomor_telepon}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="email">
								<span class="label-text font-medium">Email</span>
							</label>
							<input
								id="email"
								type="email"
								class="input-bordered input min-h-11 w-full"
								placeholder="email@contoh.com"
								bind:value={formData.email}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="website">
								<span class="label-text font-medium">Website</span>
							</label>
							<input
								id="website"
								type="url"
								class="input-bordered input min-h-11 w-full"
								placeholder="https://contoh.com"
								bind:value={formData.website}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="nomor_badan_hukum">
								<span class="label-text font-medium">Nomor Badan Hukum</span>
							</label>
							<input
								required
								id="nomor_badan_hukum"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="BH-12345-2023"
								bind:value={formData.nomor_badan_hukum}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="tahun_berdiri">
								<span class="label-text font-medium">Tahun Berdiri</span>
							</label>
							<input
								required
								id="tahun_berdiri"
								type="number"
								min="0"
								class="input-bordered input min-h-11 w-full"
								placeholder="1995"
								bind:value={formData.tahun_berdiri}
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Lokasi -->
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Lokasi</h2>
					<div class="grid grid-cols-1 gap-3 sm:gap-4">
						<div class="form-control">
							<label class="label" for="provinsi">
								<span class="label-text font-medium">Provinsi</span>
							</label>
							<input
								id="provinsi"
								type="text"
								class="input-bordered input min-h-11 w-full bg-base-200"
								value={formData.provinsi}
								readonly
							/>
						</div>
						<div class="form-control">
							<label class="label" for="kabupaten_kota">
								<span class="label-text font-medium">Kabupaten/Kota</span>
							</label>
							<input
								id="kabupaten_kota"
								type="text"
								class="input-bordered input min-h-11 w-full bg-base-200"
								value={formData.kabupaten_kota}
								readonly
							/>
						</div>
						<div class="form-control">
							<label class="label" for="kecamatan">
								<span class="label-text font-medium">Kecamatan</span>
							</label>
							<input
								id="kecamatan"
								type="text"
								class="input-bordered input min-h-11 w-full bg-base-200"
								value={formData.kecamatan}
								readonly
							/>
						</div>
						<div class="form-control">
							<label class="label" for="desa">
								<span class="label-text font-medium">Desa/Kelurahan</span>
							</label>
							<input
								required
								id="desa"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="Ciemas"
								bind:value={formData.desa}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="alamat_lengkap">
								<span class="label-text font-medium">Alamat Lengkap</span>
							</label>
							<textarea
								id="alamat_lengkap"
								class="textarea-bordered textarea min-h-20 w-full"
								placeholder="Jl. Raya Ciemas No. 45, RT 01 RW 02"
								bind:value={formData.alamat_lengkap}
								rows="2"
							></textarea>
						</div>
					</div>

					<!-- Geolocation & Map -->
					<div class="mt-4 sm:mt-6">
						<div class="mb-3 flex flex-col gap-2 sm:mb-4 sm:flex-row">
							<button
								type="button"
								class="btn min-h-11 flex-1 btn-primary sm:flex-none"
								onclick={getCurrentLocation}
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
								Lokasi Saat Ini
							</button>
						</div>

						<div class="mb-3 grid grid-cols-2 gap-2 sm:mb-4 sm:gap-4">
							<div class="form-control">
								<label class="label" for="latitude">
									<span class="label-text font-medium">Latitude</span>
								</label>
								<input
									id="latitude"
									type="number"
									step="any"
									class="input-bordered input min-h-11 w-full"
									placeholder="-6.2088"
									bind:value={formData.latitude}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="longitude">
									<span class="label-text font-medium">Longitude</span>
								</label>
								<input
									id="longitude"
									type="number"
									step="any"
									class="input-bordered input min-h-11 w-full"
									placeholder="106.8456"
									bind:value={formData.longitude}
								/>
							</div>
						</div>

						<!-- Map Container -->
						<div
							class="h-64 w-full overflow-hidden rounded-lg border border-base-300 sm:h-80 md:h-96"
						>
							<div bind:this={mapContainer} class="h-full w-full"></div>
						</div>
						<p class="mt-2 text-xs text-base-content/70 sm:text-sm">
							Klik pada peta atau geser marker untuk mengubah lokasi
						</p>
					</div>
				</div>
			</div>

			<!-- Kapasitas & Legalitas -->
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Kapasitas & Legalitas</h2>
					<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4">
						<div class="form-control">
							<label class="label" for="luas_tanah">
								<span class="label-text font-medium">Luas Tanah (m²)</span>
							</label>
							<input
								id="luas_tanah"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full"
								placeholder="1200.50"
								bind:value={formData.luas_tanah}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="luas_bangunan">
								<span class="label-text font-medium">Luas Bangunan (m²)</span>
							</label>
							<input
								id="luas_bangunan"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full"
								placeholder="800.00"
								bind:value={formData.luas_bangunan}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="status_tanah">
								<span class="label-text font-medium">Status Tanah</span>
							</label>
							<select
								id="status_tanah"
								class="select-bordered select min-h-11 w-full"
								bind:value={formData.status_tanah}
							>
								<option value="">Pilih status...</option>
								<option value="Wakaf">Wakaf</option>
								<option value="Milik Sendiri">Milik Sendiri</option>
								<option value="Sewa">Sewa</option>
								<option value="Hibah">Hibah</option>
								<option value="Lainnya">Lainnya</option>
							</select>
						</div>
						<div class="form-control">
							<label class="label" for="daya_tampung_max">
								<span class="label-text font-medium">Daya Tampung Max</span>
							</label>
							<input
								id="daya_tampung_max"
								type="number"
								min="0"
								class="input-bordered input min-h-11 w-full"
								placeholder="1000"
								bind:value={formData.daya_tampung_max}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="nomor_aiw">
								<span class="label-text font-medium">Nomor AIW</span>
							</label>
							<input
								id="nomor_aiw"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="AIW-112233"
								bind:value={formData.nomor_aiw}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="nomor_sertifikat_wakaf">
								<span class="label-text font-medium">No. Sertifikat Wakaf</span>
							</label>
							<input
								id="nomor_sertifikat_wakaf"
								type="text"
								class="input-bordered input min-h-11 w-full"
								placeholder="WKF-998877"
								bind:value={formData.nomor_sertifikat_wakaf}
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Detail Fields (Masjid Only) -->
			{#if isMasjid}
				<!-- Perpustakaan -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">Perpustakaan</h2>
						<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-4">
							<div class="form-control">
								<label class="label" for="perpustakaan_luas_m2">
									<span class="label-text font-medium">Luas (m²)</span>
								</label>
								<input
									id="perpustakaan_luas_m2"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.perpustakaan?.luas_m2 ?? ''}
									oninput={(e) => updatePerpustakaan('luas_m2', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="perpustakaan_jumlah_pengurus">
									<span class="label-text font-medium">Jml. Pengurus</span>
								</label>
								<input
									id="perpustakaan_jumlah_pengurus"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.perpustakaan?.jumlah_pengurus ?? ''}
									oninput={(e) =>
										updatePerpustakaan('jumlah_pengurus', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control col-span-2 sm:col-span-1">
								<label class="label" for="perpustakaan_kondisi">
									<span class="label-text font-medium">Kondisi</span>
								</label>
								<select
									id="perpustakaan_kondisi"
									class="select-bordered select min-h-11 w-full"
									value={formData.detail.perpustakaan?.kondisi ?? ''}
									oninput={(e) => updatePerpustakaan('kondisi', e.target.value)}
								>
									<option value="">Pilih...</option>
									<option value="Baik">Baik</option>
									<option value="Cukup Baik">Cukup Baik</option>
									<option value="Kurang Baik">Kurang Baik</option>
									<option value="Rusak">Rusak</option>
								</select>
							</div>
						</div>

						<!-- Jenis Buku List -->
						<div class="mt-3 sm:mt-4">
							<label class="label" for="jenis_buku_input">
								<span class="label-text font-medium">Jenis Buku</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="jenis_buku_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah jenis buku..."
									bind:value={newJenisBuku}
									onkeydown={(e) => handleEnterKey(e, addJenisBuku)}
								/>
								<button type="button" class="btn min-h-11 btn-secondary" onclick={addJenisBuku}>
									Add
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.perpustakaan?.jenis_buku || [] as buku, index (index)}
									<div class="badge min-h-8 gap-2 badge-primary">
										{buku}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() =>
												removeFromArray(formData.detail.perpustakaan.jenis_buku, index)}
										>
											✕
										</button>
									</div>
								{/each}
							</div>
						</div>
					</div>
				</div>

				<!-- Kalibrasi Arah Kiblat -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">Kalibrasi Kiblat</h2>
						<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4">
							<div class="form-control">
								<label class="label" for="kalibrasi_azimut">
									<span class="label-text font-medium">Azimut</span>
								</label>
								<input
									id="kalibrasi_azimut"
									type="text"
									class="input-bordered input min-h-11 w-full"
									placeholder="295.12"
									value={formData.detail.kalibrasi_arah_kiblat?.azimut ?? ''}
									oninput={(e) => updateKalibrasi('azimut', e.target.value)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="kalibrasi_tanggal">
									<span class="label-text font-medium">Tanggal</span>
								</label>
								<input
									id="kalibrasi_tanggal"
									type="date"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.kalibrasi_arah_kiblat?.tanggal_kalibrasi ?? ''}
									oninput={(e) => updateKalibrasi('tanggal_kalibrasi', e.target.value)}
								/>
							</div>
						</div>
					</div>
				</div>

				<!-- SDM Masjid -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">SDM Masjid</h2>
						<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-4">
							<div class="form-control">
								<label class="label" for="sdm_jumlah_pengurus">
									<span class="label-text font-medium">Pengurus</span>
								</label>
								<input
									id="sdm_jumlah_pengurus"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_pengurus ?? ''}
									oninput={(e) => updateSdm('jumlah_pengurus', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="sdm_jumlah_imam">
									<span class="label-text font-medium">Imam</span>
								</label>
								<input
									id="sdm_jumlah_imam"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_imam ?? ''}
									oninput={(e) => updateSdm('jumlah_imam', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="sdm_jumlah_khotib">
									<span class="label-text font-medium">Khotib</span>
								</label>
								<input
									id="sdm_jumlah_khotib"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_khotib ?? ''}
									oninput={(e) => updateSdm('jumlah_khotib', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="sdm_jumlah_muadzin">
									<span class="label-text font-medium">Muadzin</span>
								</label>
								<input
									id="sdm_jumlah_muadzin"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_muadzin ?? ''}
									oninput={(e) => updateSdm('jumlah_muadzin', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="sdm_jumlah_remaja">
									<span class="label-text font-medium">Remaja</span>
								</label>
								<input
									id="sdm_jumlah_remaja"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_remaja ?? ''}
									oninput={(e) => updateSdm('jumlah_remaja', e.target.valueAsNumber || null)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="sdm_jumlah_jemaah">
									<span class="label-text font-medium">Jemaah</span>
								</label>
								<input
									id="sdm_jumlah_jemaah"
									type="number"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.sdm_masjid?.jumlah_jemaah ?? ''}
									oninput={(e) => updateSdm('jumlah_jemaah', e.target.valueAsNumber || null)}
								/>
							</div>
						</div>
					</div>
				</div>

				<!-- Pengurus DKM -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">Pengurus DKM</h2>
						<div class="mb-3 grid grid-cols-1 gap-3 sm:mb-4 sm:grid-cols-3 sm:gap-4">
							<div class="form-control">
								<label class="label" for="pengurus_ketua">
									<span class="label-text font-medium">Ketua</span>
								</label>
								<input
									id="pengurus_ketua"
									type="text"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.pengurus_dkm?.ketua ?? ''}
									oninput={(e) => updatePengurus('ketua', e.target.value)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="pengurus_sekretaris">
									<span class="label-text font-medium">Sekretaris</span>
								</label>
								<input
									id="pengurus_sekretaris"
									type="text"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.pengurus_dkm?.sekretaris ?? ''}
									oninput={(e) => updatePengurus('sekretaris', e.target.value)}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="pengurus_bendahara">
									<span class="label-text font-medium">Bendahara</span>
								</label>
								<input
									id="pengurus_bendahara"
									type="text"
									class="input-bordered input min-h-11 w-full"
									value={formData.detail.pengurus_dkm?.bendahara ?? ''}
									oninput={(e) => updatePengurus('bendahara', e.target.value)}
								/>
							</div>
						</div>

						<!-- Nama Imam List -->
						<div class="mb-3 sm:mb-4">
							<label class="label" for="nama_imam_input">
								<span class="label-text font-medium">Nama Imam</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="nama_imam_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah nama imam..."
									bind:value={newNamaImam}
									onkeydown={(e) => handleEnterKey(e, addNamaImam)}
								/>
								<button type="button" class="btn min-h-11 btn-secondary" onclick={addNamaImam}>
									Tambah
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.pengurus_dkm?.nama_imam || [] as imam, index (index)}
									<div class="badge min-h-8 gap-2 badge-primary">
										{imam}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() => removeFromArray(formData.detail.pengurus_dkm.nama_imam, index)}
										>
											✕
										</button>
									</div>
								{/each}
							</div>
						</div>

						<!-- Nama Muazdin List -->
						<div class="mb-3 sm:mb-4">
							<label class="label" for="nama_muazdin_input">
								<span class="label-text font-medium">Nama Muazdin</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="nama_muazdin_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah nama muazdin..."
									bind:value={newNamaMuazdin}
									onkeydown={(e) => handleEnterKey(e, addNamaMuazdin)}
								/>
								<button type="button" class="btn min-h-11 btn-secondary" onclick={addNamaMuazdin}>
									Tambah
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.pengurus_dkm?.nama_muazdin || [] as muazdin, index (index)}
									<div class="badge min-h-8 gap-2 badge-primary">
										{muazdin}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() =>
												removeFromArray(formData.detail.pengurus_dkm.nama_muazdin, index)}
										>
											✕
										</button>
									</div>
								{/each}
							</div>
						</div>
					</div>
				</div>

				<!-- Fasilitas -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">Fasilitas</h2>

						<!-- Fasilitas Umum -->
						<div class="mb-3 sm:mb-4">
							<label class="label" for="fasilitas_umum_input">
								<span class="label-text font-medium">Fasilitas Umum</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="fasilitas_umum_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah fasilitas..."
									bind:value={newFasilitasUmum}
									onkeydown={(e) => handleEnterKey(e, addFasilitasUmum)}
								/>
								<button type="button" class="btn min-h-11 btn-secondary" onclick={addFasilitasUmum}>
									Tambah
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.fasilitas_umum || [] as fasilitas, index (index)}
									<div class="badge min-h-8 gap-2 badge-secondary">
										{fasilitas}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() => removeFromArray(formData.detail.fasilitas_umum, index)}
										>
											✕
										</button>
									</div>
								{/each}
							</div>
						</div>

						<!-- Fasilitas Ramah Anak -->
						<div class="mb-3 sm:mb-4">
							<label class="label" for="fasilitas_raman_anak_input">
								<span class="label-text font-medium">Fasilitas Anak</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="fasilitas_raman_anak_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah fasilitas..."
									bind:value={newFasilitasRamanAnak}
									onkeydown={(e) => handleEnterKey(e, addFasilitasRamanAnak)}
								/>
								<button
									type="button"
									class="btn min-h-11 btn-secondary"
									onclick={addFasilitasRamanAnak}
								>
									Tambah
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.fasilitas_raman_anak || [] as fasilitas, index (index)}
									<div class="badge min-h-8 gap-2 badge-accent">
										{fasilitas}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() => removeFromArray(formData.detail.fasilitas_raman_anak, index)}
										>
											✕
										</button>
									</div>
								{/each}
							</div>
						</div>

						<!-- Fasilitas Disabilitas -->
						<div class="mb-3 sm:mb-4">
							<label class="label" for="fasilitas_disabilitas_input">
								<span class="label-text font-medium">Fasilitas Disabilitas</span>
							</label>
							<div class="mb-2 flex flex-col gap-2 sm:flex-row">
								<input
									id="fasilitas_disabilitas_input"
									type="text"
									class="input-bordered input min-h-11 flex-1"
									placeholder="Tambah fasilitas..."
									bind:value={newFasilitasDisabilitas}
									onkeydown={(e) => handleEnterKey(e, addFasilitasDisabilitas)}
								/>
								<button
									type="button"
									class="btn min-h-11 btn-secondary"
									onclick={addFasilitasDisabilitas}
								>
									Tambah
								</button>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each formData.detail.fasilitas_disabilitas || [] as fasilitas, index (index)}
									<div class="badge min-h-8 gap-2 badge-outline">
										{fasilitas}
										<button
											type="button"
											class="btn btn-ghost btn-xs"
											onclick={() => removeFromArray(formData.detail.fasilitas_disabilitas, index)}
										>
											✕
										</button>
									</div>
								{/each}
								{#if formData.detail.fasilitas_disabilitas?.length === 0}
									<span class="text-sm text-base-content/50">Belum ada fasilitas disabilitas</span>
								{/if}
							</div>
						</div>
					</div>
				</div>

				<!-- Kegiatan -->
				<div class="card bg-base-100 shadow-md sm:shadow-xl">
					<div class="card-body p-4 sm:p-5">
						<h2 class="mb-3 card-title text-base sm:text-lg">Kegiatan</h2>
						<div class="mb-2 flex flex-col gap-2 sm:flex-row">
							<label class="label" for="kegiatan_input">
								<span class="label-text font-medium">Tambah Kegiatan</span>
							</label>
							<input
								id="kegiatan_input"
								type="text"
								class="input-bordered input min-h-11 flex-1"
								placeholder="Tambah kegiatan..."
								bind:value={newKegiatan}
								onkeydown={(e) => handleEnterKey(e, addKegiatan)}
							/>
							<button type="button" class="btn min-h-11 btn-secondary" onclick={addKegiatan}>
								Tambah
							</button>
						</div>
						<div class="flex flex-wrap gap-2">
							{#each formData.detail.kegiatan || [] as kegiatan, index (index)}
								<div class="badge min-h-8 gap-2 badge-success">
									{kegiatan}
									<button
										type="button"
										class="btn btn-ghost btn-xs"
										onclick={() => removeFromArray(formData.detail.kegiatan, index)}
									>
										✕
									</button>
								</div>
							{/each}
						</div>
					</div>
				</div>
			{/if}

			<!-- Submit Button -->
			<div class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8">
				<button
					onclick={() => goto('/dashboard/situs')}
					type="button"
					class="btn min-h-11 w-full btn-ghost sm:w-auto"
				>
					Batal
				</button>
				<button type="submit" class="btn min-h-11 w-full btn-primary sm:w-auto">
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
							d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
						/>
					</svg>
					Simpan Data
				</button>
			</div>
		</form>
	</div>

	<!-- Toast Notification -->
	{#if showToast}
		<div class="toast toast-end toast-top z-50">
			<div class="alert {toastType === 'success' ? 'alert-success' : 'alert-error'}">
				{#if toastType === 'success'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6 shrink-0 stroke-current"
						fill="none"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				{:else}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6 shrink-0 stroke-current"
						fill="none"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				{/if}
				<div>
					<h3 class="font-bold">{toastType === 'success' ? 'Berhasil!' : 'Gagal!'}</h3>
					<div class="text-xs">{toastMessage}</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	:global(.leaflet-container) {
		z-index: 0;
	}
</style>
