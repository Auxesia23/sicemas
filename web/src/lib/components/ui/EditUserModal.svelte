<script>
	let {
		show = false,
		isEditMode = false,
		isLoading = false,
		formData = $bindable(),
		roles = [],
		modalError = '',
		onSubmit,
		onCancel
	} = $props();
</script>

{#if show}
	<dialog class="modal" open>
		<div class="modal-box max-w-2xl">
			<h3 class="mb-6 border-b pb-2 text-xl font-bold">
				{#if isEditMode}Edit Data Petugas{:else}Registrasi Petugas Baru{/if}
			</h3>

			{#if modalError}
				<div class="mb-6 alert py-2 text-sm alert-error">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-4 w-4 shrink-0 stroke-current"
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
					<span>{modalError}</span>
				</div>
			{/if}

			<form
				onsubmit={(e) => {
					e.preventDefault();
					onSubmit();
				}}
			>
				<div class="grid grid-cols-1 gap-x-6 gap-y-4 md:grid-cols-2">
					<div class="form-control">
						<label class="label pt-0" for="nip">
							<span class="label-text font-semibold">NIP</span>
						</label>
						<input
							id="nip"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.nip}
							placeholder="18 digit angka"
							maxlength="18"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="nama">
							<span class="label-text font-semibold">Nama Lengkap</span>
						</label>
						<input
							id="nama"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.nama_lengkap}
							placeholder="Nama sesuai SK"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="jabatan">
							<span class="label-text font-semibold">Jabatan</span>
						</label>
						<input
							id="jabatan"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.jabatan}
							placeholder="Contoh: Penghulu"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="unit_kerja">
							<span class="label-text font-semibold"
								>Unit Kerja <span class="text-xs font-normal opacity-50">(Opsional)</span></span
							>
						</label>
						<input
							id="unit_kerja"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.unit_kerja}
							placeholder="Default: KUA Ciemas"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="peran">
							<span class="label-text font-semibold">Peran</span>
						</label>
						<select id="peran" class="input-bordered input w-full" bind:value={formData.peran}>
							<option value="" disabled>Pilih peran</option>
							{#each roles as role (role.ID)}
								<option value={role.name}>{role.name}</option>
							{/each}
						</select>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="email">
							<span class="label-text font-semibold">Email</span>
						</label>
						<input
							id="email"
							type="email"
							class="input-bordered input w-full"
							bind:value={formData.email}
							placeholder="petugas@kemenag.go.id"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="telp">
							<span class="label-text font-semibold">Nomor Telepon</span>
						</label>
						<input
							id="telp"
							type="tel"
							class="input-bordered input w-full"
							bind:value={formData.nomor_telepon}
							placeholder="08xxxxxxxxxx"
						/>
					</div>
				</div>

				<div class="modal-action mt-8">
					<button type="button" class="btn btn-ghost" onclick={onCancel} disabled={isLoading}>
						Batal
					</button>
					<button type="submit" class="btn px-8 btn-primary" disabled={isLoading}>
						{#if isLoading}
							<span class="loading loading-sm loading-spinner"></span>
						{/if}
						{#if isEditMode}Perbarui Data{:else}Simpan Data{/if}
					</button>
				</div>
			</form>
		</div>
		<button
			aria-label="Close modal"
			class="modal-backdrop bg-black/40 backdrop-blur-sm"
			onclick={onCancel}
		></button>
	</dialog>
{/if}
