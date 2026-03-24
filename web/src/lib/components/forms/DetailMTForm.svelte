<script>
	// Tambahkan prop errors untuk menangkap pesan validasi dari parent
	let { detail = $bindable(), errors = $bindable({}) } = $props();

	// Local state for dynamic list inputs
	let newSekretaris = $state('');
	let newBendahara = $state('');
	let newPenceramahNama = $state('');
	let newPenceramahMateri = $state('');

	// Initialize detail object if not exists
	function initDetail() {
		detail = {
			pengurus: {
				ketua: '',
				sekretaris: [],
				bendahara: []
			},
			kehadiran_pria: {
				jumlah_dewasa: 0,
				jumlah_remaja: 0,
				waktu_pengajian: ''
			},
			kehadiran_wanita: {
				jumlah_dewasa: 0,
				jumlah_remaja: 0,
				waktu_pengajian: ''
			},
			penceramah: []
		};
	}

	// Ensure detail is initialized
	if (!detail || !detail.pengurus) {
		initDetail();
	}

	// Backward compatibility jika data lama sekretaris/bendahara masih string
	if (typeof detail.pengurus?.sekretaris === 'string')
		detail.pengurus.sekretaris = [detail.pengurus.sekretaris].filter(Boolean);
	if (typeof detail.pengurus?.bendahara === 'string')
		detail.pengurus.bendahara = [detail.pengurus.bendahara].filter(Boolean);

	// Helper Functions for Arrays
	function addToArray(array, value, clearInput = true) {
		if (value.trim() && !array.includes(value.trim())) {
			array.push(value.trim());
			if (clearInput) return '';
		}
		return value;
	}

	function removeFromArray(array, index) {
		array.splice(index, 1);
	}

	function addSekretaris() {
		newSekretaris = addToArray(detail.pengurus?.sekretaris, newSekretaris);
	}
	function addBendahara() {
		newBendahara = addToArray(detail.pengurus?.bendahara, newBendahara);
	}

	// Add penceramah to the list
	function addPenceramah() {
		if (newPenceramahNama.trim() && newPenceramahMateri.trim()) {
			if (!detail.penceramah) detail.penceramah = [];
			detail.penceramah.push({
				nama: newPenceramahNama.trim(),
				materi: newPenceramahMateri.trim()
			});
			newPenceramahNama = '';
			newPenceramahMateri = '';
			errors.penceramah = null;
		}
	}

	// Remove penceramah from the list
	function removePenceramah(index) {
		detail.penceramah.splice(index, 1);
	}

	// Handle enter key for inputs
	function handleEnterKey(event, actionFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			actionFn();
		}
	}
</script>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Pengurus <span class="text-error">*</span></h2>

		<div class="mb-3">
			<label class="label" for="pengurus_ketua">
				<span class="label-text font-medium">Ketua <span class="text-error">*</span></span>
			</label>
			<input
				id="pengurus_ketua"
				type="text"
				class="input-bordered input min-h-11 w-full {errors?.pengurus?.ketua ? 'input-error' : ''}"
				bind:value={detail.pengurus.ketua}
			/>
			{#if errors?.pengurus?.ketua}
				<span class="mt-1 text-xs text-error">{errors.pengurus.ketua}</span>
			{/if}
		</div>

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			<div class="form-control">
				<label class="label" for="sekretaris_input">
					<span class="label-text font-medium">Sekretaris <span class="text-error">*</span></span>
				</label>
				<div class="flex gap-2">
					<input
						id="sekretaris_input"
						type="text"
						class="input-bordered input min-h-11 flex-1 {errors?.pengurus?.sekretaris
							? 'input-error'
							: ''}"
						placeholder="Nama Sekretaris"
						bind:value={newSekretaris}
						onkeydown={(e) => handleEnterKey(e, addSekretaris)}
					/>
					<button type="button" class="btn min-h-11 btn-secondary" onclick={addSekretaris}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus?.sekretaris || [] as person, index (index)}
						<div class="badge min-h-8 gap-2 text-white badge-info">
							{person}
							<button
								type="button"
								class="btn btn-ghost btn-xs"
								onclick={() => removeFromArray(detail.pengurus.sekretaris, index)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus?.sekretaris}
					<span class="mt-1 text-xs text-error">{errors.pengurus.sekretaris}</span>
				{/if}
			</div>

			<div class="form-control">
				<label class="label" for="bendahara_input">
					<span class="label-text font-medium">Bendahara <span class="text-error">*</span></span>
				</label>
				<div class="flex gap-2">
					<input
						id="bendahara_input"
						type="text"
						class="input-bordered input min-h-11 flex-1 {errors?.pengurus?.bendahara
							? 'input-error'
							: ''}"
						placeholder="Nama Bendahara"
						bind:value={newBendahara}
						onkeydown={(e) => handleEnterKey(e, addBendahara)}
					/>
					<button type="button" class="btn min-h-11 btn-secondary" onclick={addBendahara}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus?.bendahara || [] as person, index (index)}
						<div class="badge min-h-8 gap-2 badge-warning">
							{person}
							<button
								type="button"
								class="btn btn-ghost btn-xs"
								onclick={() => removeFromArray(detail.pengurus.bendahara, index)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus?.bendahara}
					<span class="mt-1 text-xs text-error">{errors.pengurus.bendahara}</span>
				{/if}
			</div>
		</div>
	</div>
</div>

{#if errors?.kehadiran}
	<div class="mt-4 px-2">
		<span class="text-sm font-medium text-error">{errors.kehadiran}</span>
	</div>
{/if}

<div
	class="card mt-2 bg-base-100 shadow-md sm:shadow-xl {errors?.kehadiran
		? 'border-2 border-error'
		: ''}"
>
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kehadiran Pria</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="pria_dewasa">
					<span class="label-text font-medium">Jml. Dewasa</span>
				</label>
				<input
					id="pria_dewasa"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					bind:value={detail.kehadiran_pria.jumlah_dewasa}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pria_remaja">
					<span class="label-text font-medium">Jml. Remaja</span>
				</label>
				<input
					id="pria_remaja"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					bind:value={detail.kehadiran_pria.jumlah_remaja}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pria_waktu">
					<span class="label-text font-medium">Waktu Pengajian</span>
				</label>
				<select
					id="pria_waktu"
					class="select-bordered select min-h-11 w-full"
					bind:value={detail.kehadiran_pria.waktu_pengajian}
				>
					<option value="">Pilih...</option>
					<option value="Pagi">Pagi</option>
					<option value="Sore">Sore</option>
					<option value="Malam">Malam</option>
				</select>
			</div>
		</div>
	</div>
</div>

<div
	class="card mt-4 bg-base-100 shadow-md sm:shadow-xl {errors?.kehadiran
		? 'border-2 border-error'
		: ''}"
>
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kehadiran Wanita</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="wanita_dewasa">
					<span class="label-text font-medium">Jml. Dewasa</span>
				</label>
				<input
					id="wanita_dewasa"
					type="number"
					class="input-bordered input min-h-11 w-full"
					bind:value={detail.kehadiran_wanita.jumlah_dewasa}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="wanita_remaja">
					<span class="label-text font-medium">Jml. Remaja</span>
				</label>
				<input
					id="wanita_remaja"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					bind:value={detail.kehadiran_wanita.jumlah_remaja}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="wanita_waktu">
					<span class="label-text font-medium">Waktu Pengajian</span>
				</label>
				<select
					id="wanita_waktu"
					class="select-bordered select min-h-11 w-full"
					bind:value={detail.kehadiran_wanita.waktu_pengajian}
				>
					<option value="">Pilih...</option>
					<option value="Pagi">Pagi</option>
					<option value="Sore">Sore</option>
				</select>
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">
			Daftar Penceramah <span class="text-error">*</span>
		</h2>

		<div class="mb-3 sm:mb-4">
			<label class="label" for="penceramah_inputs">
				<span class="label-text font-medium">Tambah Penceramah</span>
			</label>
			<div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-12">
				<div class="form-control sm:col-span-5">
					<input
						id="penceramah_nama"
						type="text"
						class="input-bordered input min-h-11 w-full {errors?.penceramah ? 'input-error' : ''}"
						placeholder="Nama penceramah..."
						bind:value={newPenceramahNama}
						onkeydown={(e) => handleEnterKey(e, addPenceramah)}
						oninput={() => (errors.penceramah = null)}
					/>
				</div>
				<div class="form-control sm:col-span-5">
					<input
						id="penceramah_materi"
						type="text"
						class="input-bordered input min-h-11 w-full {errors?.penceramah ? 'input-error' : ''}"
						placeholder="Materi..."
						bind:value={newPenceramahMateri}
						onkeydown={(e) => handleEnterKey(e, addPenceramah)}
						oninput={() => (errors.penceramah = null)}
					/>
				</div>
				<button
					type="button"
					class="btn min-h-11 btn-secondary sm:col-span-2"
					onclick={addPenceramah}
				>
					Tambah
				</button>
			</div>
			{#if errors?.penceramah}
				<span class="mt-1 block text-xs text-error">{errors.penceramah}</span>
			{/if}

			<div class="mt-4 flex flex-col gap-2">
				{#each detail.penceramah || [] as penceramah, index (index)}
					<div
						class="flex items-center justify-between rounded-lg border border-base-300 bg-base-200 p-3"
					>
						<div class="flex-1">
							<div class="font-medium">{penceramah.nama}</div>
							<div class="text-sm text-base-content/70">{penceramah.materi}</div>
						</div>
						<button
							type="button"
							class="btn btn-circle btn-xs btn-error"
							onclick={() => removePenceramah(index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.penceramah?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada penceramah yang ditambahkan</span>
				{/if}
			</div>
		</div>
	</div>
</div>
