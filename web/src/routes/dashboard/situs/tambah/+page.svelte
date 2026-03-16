<script>
	import { onMount, onDestroy } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import apiService from '$lib/api';
	import { goto } from '$app/navigation';
	import DetailMasjidForm from '$lib/components/forms/DetailMasjidForm.svelte';
	import DetailPesantrenForm from '$lib/components/forms/DetailPesantrenForm.svelte';
	import DetailMusholaForm from '$lib/components/forms/DetailMusholaForm.svelte';
	import DetailMTForm from '$lib/components/forms/DetailMTForm.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';

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

	// Tipologi options based on jenis situs
	const tipologiOptions = {
		Masjid: ['Masjid Besar', 'Masjid Jami', 'Masjid Bersejarah', 'Masjid Publik'],
		Mushola: ['Mushola'],
		Ponpes: ['Salafiyah', 'Kholafiyah', 'Salafiyah Wajar Dikdas', 'Muadalah'],
		MT: ['Majelis Taklim'],
		Madrasah: ['Madrasah Diniyah', 'Madrasah Ibtidaiyah', 'Madrasah Tsanawiyah', 'Madrasah Aliyah'],
		Lainnya: ['Umum', 'Khusus']
	};

	// Get tipologi options for selected jenis situs
	function getTipologiOptions() {
		if (!selectedJenisSitusName) return [];
		// Find matching key in tipologiOptions
		for (const [key, options] of Object.entries(tipologiOptions)) {
			if (selectedJenisSitusName.toLowerCase().includes(key.toLowerCase())) {
				return options;
			}
		}
		// Default options if no match
		return ['Umum', 'Khusus'];
	}

	// State for site types dropdown
	let jenisSitusList = $state([]);
	let selectedJenisSitusName = $state('');

	// State for map
	let mapContainer;
	let map;
	let marker;

	// State for toast notification
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success'); // 'success' or 'error'

	// State for images
	let selectedImages = $state([]); // Array of { file, preview, webpBlob }
	let isConverting = $state(false);
	let convertedImages = $state([]); // Array of webp blobs ready to upload
	let fileInputRef; // Reference to file input element
	let isSubmitting = $state(false); // Loading state for submit button

	// Derived: check if selected site is Masjid
	let isMasjid = $derived(selectedJenisSitusName.toLowerCase().includes('masjid'));

	// Derived: check if selected site is Pesantren
	let isPesantren = $derived(selectedJenisSitusName.toLowerCase().includes('ponpes'));

	// Derived: check if selected site is Musholla
	let isMusholla = $derived(selectedJenisSitusName.toLowerCase().includes('mushola'));

	// Derived: check if selected site is Majelis Taklim
	let isMT = $derived(selectedJenisSitusName.toLowerCase().includes('mt'));

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

	// Cleanup when component unmounts
	onDestroy(() => {
		// Clean up image preview URLs
		selectedImages.forEach((img) => {
			if (img?.preview) {
				URL.revokeObjectURL(img.preview);
			}
		});
		// Reset file input
		if (fileInputRef) {
			fileInputRef.value = '';
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
			formData.jenis_tipologi = '';

			// Initialize or reset detail based on selected site type
			const selectedName = selected.nama || selected.name || '';
			if (selectedName.toLowerCase().includes('masjid')) {
				// Masjid detail will be initialized by child component
				formData.detail = {};
				// Auto-set tipologi to first option
				formData.jenis_tipologi = tipologiOptions.Masjid[0];
			} else if (selectedName.toLowerCase().includes('ponpes')) {
				// Pesantren detail will be initialized by child component
				formData.detail = {};
				// Auto-set tipologi to first option
				formData.jenis_tipologi = tipologiOptions.Ponpes[0];
			} else if (selectedName.toLowerCase().includes('mushola')) {
				// Musholla detail will be initialized by child component
				formData.detail = {};
				// Auto-set tipologi to Musholla
				formData.jenis_tipologi = tipologiOptions.Mushola[0];
			} else if (selectedName.toLowerCase().includes('mt')) {
				// Majelis Taklim detail will be initialized by child component
				formData.detail = {};
				// Auto-set tipologi to Majelis Taklim
				formData.jenis_tipologi = tipologiOptions.MT[0];
			} else {
				// Clear detail for other site types
				formData.detail = {};
			}
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

	// Convert image to WebP
	async function convertToWebP(file) {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = (e) => {
				const img = new Image();
				img.onload = () => {
					const canvas = document.createElement('canvas');
					canvas.width = img.width;
					canvas.height = img.height;
					const ctx = canvas.getContext('2d');
					ctx.drawImage(img, 0, 0);
					canvas.toBlob(
						(blob) => {
							if (blob) {
								resolve(blob);
							} else {
								reject(new Error('Failed to convert image'));
							}
						},
						'image/webp',
						0.8
					);
				};
				img.onerror = () => reject(new Error('Failed to load image'));
				img.src = e.target.result;
			};
			reader.onerror = () => reject(new Error('Failed to read file'));
			reader.readAsDataURL(file);
		});
	}

	// Handle image selection
	async function handleImageSelect(event) {
		const files = Array.from(event.target.files);
		if (files.length === 0) return;

		isConverting = true;

		try {
			for (const file of files) {
				if (!file.type.startsWith('image/')) continue;

				// Create preview URL
				const preview = URL.createObjectURL(file);

				// Convert to WebP
				const webpBlob = await convertToWebP(file);

				selectedImages.push({
					file,
					preview,
					webpBlob
				});
			}
		} catch (error) {
			console.error('Failed to convert image:', error);
			toastMessage = 'Gagal mengkonversi gambar';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 3000);
		} finally {
			isConverting = false;
		}
	}

	// Remove selected image
	function removeImage(index) {
		if (selectedImages[index]?.preview) {
			URL.revokeObjectURL(selectedImages[index].preview);
		}
		selectedImages.splice(index, 1);

		// Reset file input if no more images
		if (selectedImages.length === 0 && fileInputRef) {
			fileInputRef.value = '';
		}
	}

	// Upload images to API
	async function uploadImages(situsId) {
		if (selectedImages.length === 0) return;

		const formDataImages = new FormData();
		for (const img of selectedImages) {
			formDataImages.append('images', img.webpBlob, `image-${Date.now()}.webp`);
		}

		try {
			const response = await apiService.post(`/situs/${situsId}/foto`, formDataImages);
			if (!response.ok) {
				console.error('Failed to upload images:', await response.text());
			}
		} catch (error) {
			console.error('Image upload error:', error);
			throw error;
		}
	}

	// Submit handler
	async function handleSubmit() {
		if (isSubmitting) return; // Prevent multiple submissions

		isSubmitting = true;

		try {
			const response = await apiService.post('/situs', formData);
			const data = await response.json();

			if (response.status === 201) {
				const situsId = data.id;

				// Upload images if any
				if (situsId && selectedImages.length > 0) {
					await uploadImages(situsId);
				}

				// Show success toast notification
				toastMessage = 'Data situs berhasil disimpan';
				toastType = 'success';
				showToast = true;

				// Clean up preview URLs
				selectedImages.forEach((img) => {
					if (img?.preview) {
						URL.revokeObjectURL(img.preview);
					}
				});
				selectedImages = [];

				// Reset file input
				if (fileInputRef) {
					fileInputRef.value = '';
				}

				// REDIRECT (isSubmitting tetap TRUE di sini)
				setTimeout(() => {
					showToast = false;
					goto('/dashboard/situs');
				}, 2000);
			} else {
				toastMessage = 'Gagal menyimpan data situs';
				toastType = 'error';
				showToast = true;
				setTimeout(() => {
					showToast = false;
				}, 5000);

				// Jika gagal API (bukan 201), hidupkan lagi tombolnya
				isSubmitting = false;
			}
		} catch (error) {
			console.error('Failed to submit form:', error);
			toastMessage = 'Terjadi kesalahan. Silakan coba lagi.';
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 5000);

			// Jika masuk ke catch blok karena error network, hidupkan lagi tombolnya
			isSubmitting = false;
		}
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
								bind:value={formData.nama}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="jenis_tipologi">
								<span class="label-text font-medium">Jenis Tipologi</span>
							</label>

							{#if getTipologiOptions().length <= 1}
								<input
									id="jenis_tipologi"
									type="text"
									class="input-bordered input min-h-11 w-full bg-base-200"
									value={formData.jenis_tipologi}
									readonly
								/>
								<label for="jenis_tipologi" class="label">
									<span class="label-text-alt text-base-content/70">
										Tipologi otomatis terisi
									</span>
								</label>
							{:else}
								<select
									id="jenis_tipologi"
									class="select-bordered select min-h-11 w-full"
									required
									bind:value={formData.jenis_tipologi}
								>
									{#each getTipologiOptions() as option}
										<option value={option}>{option}</option>
									{/each}
								</select>
							{/if}
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
								bind:value={formData.desa}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="alamat_lengkap">
								<span class="label-text font-medium">Alamat Lengkap</span>
							</label>
							<textarea
								required
								id="alamat_lengkap"
								class="textarea-bordered textarea min-h-20 w-full"
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
									required
									id="latitude"
									type="number"
									step="any"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.latitude}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="longitude">
									<span class="label-text font-medium">Longitude</span>
								</label>
								<input
									required
									id="longitude"
									type="number"
									step="any"
									class="input-bordered input min-h-11 w-full"
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
								required
								id="luas_tanah"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full"
								bind:value={formData.luas_tanah}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="luas_bangunan">
								<span class="label-text font-medium">Luas Bangunan (m²)</span>
							</label>
							<input
								required
								id="luas_bangunan"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full"
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
								required
								id="daya_tampung_max"
								type="number"
								min="0"
								class="input-bordered input min-h-11 w-full"
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
								bind:value={formData.nomor_sertifikat_wakaf}
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Detail Fields (Masjid Only) -->
			{#if isMasjid}
				<DetailMasjidForm bind:detail={formData.detail} />
			{/if}

			<!-- Detail Fields (Pesantren Only) -->
			{#if isPesantren}
				<DetailPesantrenForm bind:detail={formData.detail} />
			{/if}

			<!-- Detail Fields (Musholla Only) -->
			{#if isMusholla}
				<DetailMusholaForm bind:detail={formData.detail} />
			{/if}

			<!-- Detail Fields (Majelis Taklim Only) -->
			{#if isMT}
				<DetailMTForm bind:detail={formData.detail} />
			{/if}

			<!-- Upload Gambar -->
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Foto Situs</h2>
					<div class="form-control">
						<label class="label" for="images">
							<span class="label-text font-medium">Unggah Foto</span>
						</label>
						<div class="mb-2 flex flex-col gap-3">
							<input
								id="images"
								type="file"
								accept="image/*"
								multiple
								class="file-input-bordered file-input min-h-11 w-full"
								onchange={handleImageSelect}
								disabled={isConverting}
								bind:this={fileInputRef}
							/>
							{#if isConverting}
								<div class="flex items-center gap-2 text-sm text-base-content/70">
									<span class="loading loading-sm loading-spinner"></span>
									Mengkonversi gambar ke WebP...
								</div>
							{/if}
						</div>

						<!-- Image Previews -->
						{#if selectedImages.length > 0}
							<div class="mt-2 grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
								{#each selectedImages as img, index (img.preview)}
									<div class="relative overflow-hidden rounded-lg border border-base-300">
										<img src={img.preview} alt="Preview" class="h-32 w-full object-cover" />
										<button
											type="button"
											class="btn absolute top-1 right-1 btn-circle btn-xs btn-error"
											onclick={() => removeImage(index)}
										>
											✕
										</button>
										<div
											class="absolute right-1 bottom-1 rounded bg-black/70 px-1 text-xs text-white"
										>
											WebP
										</div>
									</div>
								{/each}
							</div>
							<p class="mt-2 text-xs text-base-content/70">
								{selectedImages.length} foto dipilih. Format akan dikonversi ke WebP sebelum diunggah.
							</p>
						{/if}
					</div>
				</div>
			</div>

			<!-- Submit Button -->
			<div class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8">
				<button
					onclick={() => goto('/dashboard/situs')}
					type="button"
					class="btn min-h-11 w-full btn-ghost sm:w-auto"
					disabled={isSubmitting}
				>
					Batal
				</button>
				<button
					type="submit"
					class="btn min-h-11 w-full btn-primary sm:w-auto"
					disabled={isSubmitting}
				>
					{#if isSubmitting}
						<span class="loading loading-spinner"></span>
						Menyimpan...
					{:else}
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
