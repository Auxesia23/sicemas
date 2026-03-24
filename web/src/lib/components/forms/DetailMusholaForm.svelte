<script>
	// Tambahkan prop errors untuk menangkap pesan validasi dari parent
	let { detail = $bindable(), errors = $bindable({}) } = $props();

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
				jumlah_imam: 0,
				jumlah_muadzin: 0,
				jumlah_jemaah: 0
			},
			fasilitas_umum: [],
			nama_imam: [],
			nama_muadzin: [],
			kegiatan: []
		};
	}

	// Ensure detail is initialized
	if (!detail || !detail.kalibrasi_arah_kiblat) {
		initDetail();
	}

	// Helper penambah array (Badge)
	function addToArray(array, value) {
		if (value.trim() && !array.includes(value.trim())) {
			array.push(value.trim());
			return '';
		}
		return value;
	}

	function removeFromArray(array, index) {
		array.splice(index, 1);
	}

	// Event Handlers
	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}

	// Fungsi hapus error saat mengetik
	function clearError(parent, child) {
		if (errors?.[parent]?.[child]) {
			errors[parent][child] = null;
		}
		// Jika error ada di root level (seperti nama_imam)
		if (errors?.[parent] && !child) {
			errors[parent] = null;
		}
	}
</script>

<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">Kalibrasi Kiblat</h2>
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
					bind:value={detail.kalibrasi_arah_kiblat.azimut}
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
					bind:value={detail.kalibrasi_arah_kiblat.tanggal_kalibrasi}
				/>
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			SDM Mushola <span class="text-error">*</span>
		</h2>
		{#if errors?.sdm_masjid && typeof errors.sdm_masjid === 'string'}
			<div class="mb-2 text-sm text-error">{errors.sdm_masjid}</div>
		{/if}
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			{#each [{ id: 'imam', label: 'Imam', path: 'jumlah_imam' }, { id: 'muadzin', label: 'Muadzin', path: 'jumlah_muadzin' }, { id: 'jemaah', label: 'Jemaah', path: 'jumlah_jemaah' }] as item}
				<div class="form-control">
					<label class="label" for="sdm_{item.id}">
						<span class="label-text font-medium"
							>{item.label} <span class="text-error">*</span></span
						>
					</label>
					<input
						id="sdm_{item.id}"
						type="number"
						min="0"
						class="input-bordered input min-h-11 w-full {errors?.sdm_masjid?.[item.path]
							? 'input-error'
							: ''}"
						placeholder="0"
						bind:value={detail.sdm_masjid[item.path]}
						oninput={() => clearError('sdm_masjid', item.path)}
					/>
					{#if errors?.sdm_masjid?.[item.path]}
						<span class="mt-1 text-xs text-error">{errors.sdm_masjid[item.path]}</span>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">Fasilitas Umum</h2>
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
					onkeydown={(e) =>
						handleEnterKey(
							e,
							() => (newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum))
						)}
				/>
				<button
					type="button"
					class="btn min-h-11 btn-secondary"
					onclick={() => (newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum))}
				>
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
							onclick={() => removeFromArray(detail.fasilitas_umum, index)}>✕</button
						>
					</div>
				{/each}
				{#if detail.fasilitas_umum?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada fasilitas umum</span>
				{/if}
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">Pengurus & Kegiatan</h2>

		<div class="mb-4">
			<label class="label" for="nama_imam_input">
				<span class="label-text font-medium">Nama Imam <span class="text-error">*</span></span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="nama_imam_input"
					type="text"
					class="input-bordered input min-h-11 flex-1 {errors?.nama_imam ? 'input-error' : ''}"
					placeholder="Tambah nama imam..."
					bind:value={newNamaImam}
					onkeydown={(e) =>
						handleEnterKey(e, () => {
							newNamaImam = addToArray(detail.nama_imam, newNamaImam);
							clearError('nama_imam');
						})}
				/>
				<button
					type="button"
					class="btn min-h-11 btn-secondary"
					onclick={() => {
						newNamaImam = addToArray(detail.nama_imam, newNamaImam);
						clearError('nama_imam');
					}}
				>
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
							onclick={() => removeFromArray(detail.nama_imam, index)}>✕</button
						>
					</div>
				{/each}
			</div>
			{#if errors?.nama_imam}
				<span class="mt-1 block text-xs text-error">{errors.nama_imam}</span>
			{/if}
		</div>

		<div class="mb-4">
			<label class="label" for="nama_muazdin_input">
				<span class="label-text font-medium">Nama Muadzin <span class="text-error">*</span></span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="nama_muazdin_input"
					type="text"
					class="input-bordered input min-h-11 flex-1 {errors?.nama_muadzin ? 'input-error' : ''}"
					placeholder="Tambah nama muadzin..."
					bind:value={newNamaMuazdin}
					onkeydown={(e) =>
						handleEnterKey(e, () => {
							newNamaMuazdin = addToArray(detail.nama_muadzin, newNamaMuazdin);
							clearError('nama_muadzin');
						})}
				/>
				<button
					type="button"
					class="btn min-h-11 btn-secondary"
					onclick={() => {
						newNamaMuazdin = addToArray(detail.nama_muadzin, newNamaMuazdin);
						clearError('nama_muadzin');
					}}
				>
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
							onclick={() => removeFromArray(detail.nama_muadzin, index)}>✕</button
						>
					</div>
				{/each}
			</div>
			{#if errors?.nama_muadzin}
				<span class="mt-1 block text-xs text-error">{errors.nama_muadzin}</span>
			{/if}
		</div>

		<div>
			<label class="label" for="kegiatan_input">
				<span class="label-text font-medium">Kegiatan (Opsional)</span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="kegiatan_input"
					type="text"
					class="input-bordered input min-h-11 flex-1"
					placeholder="Tambah kegiatan..."
					bind:value={newKegiatan}
					onkeydown={(e) =>
						handleEnterKey(e, () => (newKegiatan = addToArray(detail.kegiatan, newKegiatan)))}
				/>
				<button
					type="button"
					class="btn min-h-11 btn-secondary"
					onclick={() => (newKegiatan = addToArray(detail.kegiatan, newKegiatan))}
				>
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
							onclick={() => removeFromArray(detail.kegiatan, index)}>✕</button
						>
					</div>
				{/each}
				{#if detail.kegiatan?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada kegiatan</span>
				{/if}
			</div>
		</div>
	</div>
</div>
