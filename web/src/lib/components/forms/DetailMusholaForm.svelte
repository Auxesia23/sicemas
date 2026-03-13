<script>
	// Svelte 5 Runes: Bindable prop for two-way binding with parent
	let { detail = $bindable() } = $props();

	// Local state for dynamic list inputs
	let newFasilitasUmum = $state('');
	let newNamaImam = $state('');
	let newNamaMuazdin = $state('');
	let newKegiatan = $state('');

	// Initialize detail object if not exists
	function initDetail() {
		detail = {
			kalibrasi_arah_kiblat: {
				azimut: '',
				tanggal_kalibrasi: ''
			},
			sdm_masjid: {
				jumlah_imam: null,
				jumlah_muadzin: null,
				jumlah_jemaah: null
			},
			fasilitas_umum: [],
			nama_imam: [],
			nama_muadzin: [],
			kegiatan: []
		};
	}

	// Ensure detail is initialized (synchronous, NOT in $effect)
	if (!detail || !detail.kalibrasi_arah_kiblat) {
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

	function addFasilitasUmum() {
		newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum, true);
	}

	function addNamaImam() {
		newNamaImam = addToArray(detail.nama_imam, newNamaImam, true);
	}

	function addNamaMuazdin() {
		newNamaMuazdin = addToArray(detail.nama_muadzin, newNamaMuazdin, true);
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
	function updateKalibrasi(field, value) {
		if (!detail.kalibrasi_arah_kiblat) initDetail();
		detail.kalibrasi_arah_kiblat[field] = value;
	}

	function updateSdm(field, value) {
		if (!detail.sdm_masjid) initDetail();
		detail.sdm_masjid[field] = value;
	}
</script>

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

<!-- SDM Mushola -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">SDM Mushola</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="sdm_jumlah_imam">
					<span class="label-text font-medium">Imam</span>
				</label>
				<input
					id="sdm_jumlah_imam"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.sdm_masjid?.jumlah_imam ?? ''}
					oninput={(e) => updateSdm('jumlah_imam', e.target.valueAsNumber || null)}
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
					placeholder="0"
					value={detail.sdm_masjid?.jumlah_muadzin ?? ''}
					oninput={(e) => updateSdm('jumlah_muadzin', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="sdm_jumlah_jemaah">
					<span class="label-text font-medium">Jemaah</span>
				</label>
				<input
					id="sdm_jumlah_jemaah"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.sdm_masjid?.jumlah_jemaah ?? ''}
					oninput={(e) => updateSdm('jumlah_jemaah', e.target.valueAsNumber || null)}
				/>
			</div>
		</div>
	</div>
</div>

<!-- Fasilitas Umum -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Fasilitas Umum</h2>
		<div class="mb-3 sm:mb-4">
			<label class="label" for="fasilitas_umum_input">
				<span class="label-text font-medium">Tambah Fasilitas</span>
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
				{#if detail.fasilitas_umum?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada fasilitas umum</span>
				{/if}
			</div>
		</div>
	</div>
</div>

<!-- Pengurus & Kegiatan -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Pengurus & Kegiatan</h2>

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
				{#each detail.nama_imam || [] as imam, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{imam}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.nama_imam, index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.nama_imam?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada nama imam</span>
				{/if}
			</div>
		</div>

		<!-- Nama Muadzin List -->
		<div class="mb-3 sm:mb-4">
			<label class="label" for="nama_muazdin_input">
				<span class="label-text font-medium">Nama Muadzin</span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="nama_muazdin_input"
					type="text"
					class="input-bordered input min-h-11 flex-1"
					placeholder="Tambah nama muadzin..."
					bind:value={newNamaMuazdin}
					onkeydown={(e) => handleEnterKey(e, addNamaMuazdin)}
				/>
				<button type="button" class="btn min-h-11 btn-secondary" onclick={addNamaMuazdin}>
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.nama_muadzin || [] as muadzin, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{muadzin}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.nama_muadzin, index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.nama_muadzin?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada nama muadzin</span>
				{/if}
			</div>
		</div>

		<!-- Kegiatan List -->
		<div class="mb-3 sm:mb-4">
			<label class="label" for="kegiatan_input">
				<span class="label-text font-medium">Kegiatan</span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
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
				{#if detail.kegiatan?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada kegiatan</span>
				{/if}
			</div>
		</div>
	</div>
</div>
