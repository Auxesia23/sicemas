<script>
	// Svelte 5 Runes: Bindable prop for two-way binding with parent
	let { detail = $bindable() } = $props();

	// Local state for dynamic list inputs
	let newJenisBuku = $state('');
	let newNamaImam = $state('');
	let newNamaMuazdin = $state('');
	let newFasilitasUmum = $state('');
	let newFasilitasRamanAnak = $state('');
	let newFasilitasDisabilitas = $state('');
	let newKegiatan = $state('');

	// Initialize detail object if not exists
	function initDetail() {
		detail = {
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

	// Ensure detail is initialized
	if (!detail?.perpustakaan) {
		initDetail();
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
		newJenisBuku = addToArray(detail.perpustakaan?.jenis_buku, newJenisBuku, true);
	}

	function addNamaImam() {
		newNamaImam = addToArray(detail.pengurus_dkm?.nama_imam, newNamaImam, true);
	}

	function addNamaMuazdin() {
		newNamaMuazdin = addToArray(detail.pengurus_dkm?.nama_muazdin, newNamaMuazdin, true);
	}

	function addFasilitasUmum() {
		newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum, true);
	}

	function addFasilitasRamanAnak() {
		newFasilitasRamanAnak = addToArray(detail.fasilitas_raman_anak, newFasilitasRamanAnak, true);
	}

	function addFasilitasDisabilitas() {
		newFasilitasDisabilitas = addToArray(
			detail.fasilitas_disabilitas,
			newFasilitasDisabilitas,
			true
		);
	}

	function addKegiatan() {
		newKegiatan = addToArray(detail.kegiatan, newKegiatan, true);
	}

	// Handle enter key for list inputs
	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}

	// Helper functions for detail field bindings
	function updatePerpustakaan(field, value) {
		if (!detail.perpustakaan) initDetail();
		detail.perpustakaan[field] = value;
	}

	function updateKalibrasi(field, value) {
		if (!detail.kalibrasi_arah_kiblat) initDetail();
		detail.kalibrasi_arah_kiblat[field] = value;
	}

	function updateSdm(field, value) {
		if (!detail.sdm_masjid) initDetail();
		detail.sdm_masjid[field] = value;
	}

	function updatePengurus(field, value) {
		if (!detail.pengurus_dkm) initDetail();
		detail.pengurus_dkm[field] = value;
	}
</script>

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
					value={detail.perpustakaan?.luas_m2 ?? ''}
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
					value={detail.perpustakaan?.jumlah_pengurus ?? ''}
					oninput={(e) => updatePerpustakaan('jumlah_pengurus', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control col-span-2 sm:col-span-1">
				<label class="label" for="perpustakaan_kondisi">
					<span class="label-text font-medium">Kondisi</span>
				</label>
				<select
					id="perpustakaan_kondisi"
					class="select-bordered select min-h-11 w-full"
					value={detail.perpustakaan?.kondisi ?? ''}
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
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.perpustakaan?.jenis_buku || [] as buku, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{buku}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.perpustakaan.jenis_buku, index)}
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
					value={detail.kalibrasi_arah_kiblat?.azimut ?? ''}
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
					value={detail.kalibrasi_arah_kiblat?.tanggal_kalibrasi ?? ''}
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
					value={detail.sdm_masjid?.jumlah_pengurus ?? ''}
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
					value={detail.sdm_masjid?.jumlah_imam ?? ''}
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
					value={detail.sdm_masjid?.jumlah_khotib ?? ''}
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
					value={detail.sdm_masjid?.jumlah_muadzin ?? ''}
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
					value={detail.sdm_masjid?.jumlah_remaja ?? ''}
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
					value={detail.sdm_masjid?.jumlah_jemaah ?? ''}
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
					value={detail.pengurus_dkm?.ketua ?? ''}
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
					value={detail.pengurus_dkm?.sekretaris ?? ''}
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
					value={detail.pengurus_dkm?.bendahara ?? ''}
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
				{#each detail.pengurus_dkm?.nama_imam || [] as imam, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{imam}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.pengurus_dkm.nama_imam, index)}
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
				{#each detail.pengurus_dkm?.nama_muazdin || [] as muazdin, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{muazdin}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.pengurus_dkm.nama_muazdin, index)}
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
				{#each detail.fasilitas_umum || [] as fasilitas, index (index)}
					<div class="badge min-h-8 gap-2 badge-secondary">
						{fasilitas}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.fasilitas_umum, index)}
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
				<button type="button" class="btn min-h-11 btn-secondary" onclick={addFasilitasRamanAnak}>
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_raman_anak || [] as fasilitas, index (index)}
					<div class="badge min-h-8 gap-2 badge-accent">
						{fasilitas}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.fasilitas_raman_anak, index)}
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
				<button type="button" class="btn min-h-11 btn-secondary" onclick={addFasilitasDisabilitas}>
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_disabilitas || [] as fasilitas, index (index)}
					<div class="badge min-h-8 gap-2 badge-outline">
						{fasilitas}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.fasilitas_disabilitas, index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.fasilitas_disabilitas?.length === 0}
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
			{#each detail.kegiatan || [] as kegiatan, index (index)}
				<div class="badge min-h-8 gap-2 badge-success">
					{kegiatan}
					<button
						type="button"
						class="btn btn-ghost btn-xs"
						onclick={() => removeFromArray(detail.kegiatan, index)}
					>
						✕
					</button>
				</div>
			{/each}
		</div>
	</div>
</div>
