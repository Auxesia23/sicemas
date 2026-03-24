<script>
	// Tambahkan prop errors
	let { detail = $bindable(), errors = $bindable({}) } = $props();

	let newFasilitasPondok = $state('');
	let newSekretaris = $state('');
	let newBendahara = $state('');

	function initDetail() {
		detail = {
			nama_yayasan: '',
			pimpinan_pesantren: '',
			kepengurusan: {
				ketua: '',
				sekretaris: [],
				bendahara: []
			},
			santri: {
				mondok: { pria: 0, wanita: 0 },
				tidak_mondok: { pria: 0, wanita: 0 },
				total: { pria: 0, wanita: 0 }
			},
			fasilitas_pondok: []
		};
	}

	if (!detail?.nama_yayasan) {
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

	// Fungsi hapus error saat mengetik
	function clearError(parent, sub, field) {
		if (field) {
			// Untuk level 3 (santri.mondok.pria)
			if (errors?.[parent]?.[sub]?.[field]) errors[parent][sub][field] = null;
		} else if (sub) {
			// Untuk level 2 (kepengurusan.ketua)
			if (errors?.[parent]?.[sub]) errors[parent][sub] = null;
		} else {
			// Untuk level 1 (nama_yayasan)
			if (errors?.[parent]) errors[parent] = null;
		}
	}

	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}
</script>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			Informasi Yayasan <span class="text-error">*</span>
		</h2>
		<div class="grid grid-cols-1 gap-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="nama_yayasan">
					<span class="label-text font-medium">Nama Yayasan <span class="text-error">*</span></span>
				</label>
				<input
					id="nama_yayasan"
					type="text"
					class="input-bordered input min-h-11 w-full {errors?.nama_yayasan ? 'input-error' : ''}"
					bind:value={detail.nama_yayasan}
					oninput={() => clearError('nama_yayasan')}
				/>
				{#if errors?.nama_yayasan}<span class="mt-1 text-xs text-error">{errors.nama_yayasan}</span
					>{/if}
			</div>
			<div class="form-control">
				<label class="label" for="pimpinan_pesantren">
					<span class="label-text font-medium"
						>Pimpinan Pesantren <span class="text-error">*</span></span
					>
				</label>
				<input
					id="pimpinan_pesantren"
					type="text"
					class="input-bordered input min-h-11 w-full {errors?.pimpinan_pesantren
						? 'input-error'
						: ''}"
					bind:value={detail.pimpinan_pesantren}
					oninput={() => clearError('pimpinan_pesantren')}
				/>
				{#if errors?.pimpinan_pesantren}<span class="mt-1 text-xs text-error"
						>{errors.pimpinan_pesantren}</span
					>{/if}
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			Kepengurusan <span class="text-error">*</span>
		</h2>

		<div class="form-control mb-4">
			<label class="label" for="kepengurusan_ketua">
				<span class="label-text font-medium"
					>Ketua Yayasan/Pesantren <span class="text-error">*</span></span
				>
			</label>
			<input
				id="kepengurusan_ketua"
				type="text"
				class="input-bordered input min-h-11 w-full {errors?.kepengurusan?.ketua
					? 'input-error'
					: ''}"
				bind:value={detail.kepengurusan.ketua}
				oninput={() => clearError('kepengurusan', 'ketua')}
			/>
			{#if errors?.kepengurusan?.ketua}<span class="mt-1 text-xs text-error"
					>{errors.kepengurusan.ketua}</span
				>{/if}
		</div>

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			<div class="form-control">
				<label class="label" for="sek_in"
					><span class="label-text font-medium">Sekretaris <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="sek_in"
						type="text"
						class="input-bordered input flex-1 {errors?.kepengurusan?.sekretaris
							? 'input-error'
							: ''}"
						placeholder="Tambah Sekretaris..."
						bind:value={newSekretaris}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() => (newSekretaris = addToArray(detail.kepengurusan.sekretaris, newSekretaris))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => {
							newSekretaris = addToArray(detail.kepengurusan.sekretaris, newSekretaris);
							clearError('kepengurusan', 'sekretaris');
						}}>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.kepengurusan.sekretaris as p, i}
						<div class="badge h-8 gap-2 text-white badge-info">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.kepengurusan.sekretaris, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.kepengurusan?.sekretaris}<span class="mt-1 text-xs text-error"
						>{errors.kepengurusan.sekretaris}</span
					>{/if}
			</div>

			<div class="form-control">
				<label class="label" for="ben_in"
					><span class="label-text font-medium">Bendahara <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="ben_in"
						type="text"
						class="input-bordered input flex-1 {errors?.kepengurusan?.bendahara
							? 'input-error'
							: ''}"
						placeholder="Tambah Bendahara..."
						bind:value={newBendahara}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() => (newBendahara = addToArray(detail.kepengurusan.bendahara, newBendahara))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => {
							newBendahara = addToArray(detail.kepengurusan.bendahara, newBendahara);
							clearError('kepengurusan', 'bendahara');
						}}>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.kepengurusan.bendahara as p, i}
						<div class="badge h-8 gap-2 badge-warning">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.kepengurusan.bendahara, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.kepengurusan?.bendahara}<span class="mt-1 text-xs text-error"
						>{errors.kepengurusan.bendahara}</span
					>{/if}
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			Data Santri <span class="text-error">*</span>
		</h2>

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
			{#each ['mondok', 'tidak_mondok', 'total'] as cat}
				<div class="card bg-base-200">
					<div class="card-body p-3">
						<h3 class="card-title text-sm capitalize">{cat.replace('_', ' ')}</h3>
						<div class="mt-2 grid grid-cols-2 gap-2">
							<div class="form-control">
								<label class="label p-1" for="{cat}_pria"
									><span class="label-text text-[10px]">Pria</span></label
								>
								<input
									id="{cat}_pria"
									type="number"
									min="0"
									class="input-bordered input input-sm w-full {errors?.santri?.[cat]?.pria
										? 'input-error'
										: ''}"
									bind:value={detail.santri[cat].pria}
									oninput={() => clearError('santri', cat, 'pria')}
								/>
								{#if errors?.santri?.[cat]?.pria}<span class="mt-0.5 text-[10px] text-error"
										>{errors.santri[cat].pria}</span
									>{/if}
							</div>
							<div class="form-control">
								<label class="label p-1" for="{cat}_wanita"
									><span class="label-text text-[10px]">Wanita</span></label
								>
								<input
									id="{cat}_wanita"
									type="number"
									min="0"
									class="input-bordered input input-sm w-full {errors?.santri?.[cat]?.wanita
										? 'input-error'
										: ''}"
									bind:value={detail.santri[cat].wanita}
									oninput={() => clearError('santri', cat, 'wanita')}
								/>
								{#if errors?.santri?.[cat]?.wanita}<span class="mt-0.5 text-[10px] text-error"
										>{errors.santri[cat].wanita}</span
									>{/if}
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">Fasilitas Pondok</h2>
		<div class="mb-3 sm:mb-4">
			<label class="label" for="fasilitas_pondok_input">
				<span class="label-text font-medium">Tambah Fasilitas</span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="fasilitas_pondok_input"
					type="text"
					class="input-bordered input min-h-11 flex-1"
					placeholder="Tambah fasilitas..."
					bind:value={newFasilitasPondok}
					onkeydown={(e) =>
						handleEnterKey(
							e,
							() => (newFasilitasPondok = addToArray(detail.fasilitas_pondok, newFasilitasPondok))
						)}
				/>
				<button
					type="button"
					class="btn min-h-11 btn-secondary"
					onclick={() =>
						(newFasilitasPondok = addToArray(detail.fasilitas_pondok, newFasilitasPondok))}
				>
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_pondok || [] as fasilitas, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{fasilitas}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.fasilitas_pondok, index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.fasilitas_pondok?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada fasilitas pondok</span>
				{/if}
			</div>
		</div>
	</div>
</div>
