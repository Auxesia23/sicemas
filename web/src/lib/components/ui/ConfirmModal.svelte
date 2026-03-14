<script>
	let {
		show = false,
		title = 'Konfirmasi',
		message = '',
		confirmText = 'Hapus',
		cancelText = 'Batal',
		isProcessing = false,
		onConfirm,
		onCancel
	} = $props();
</script>

{#if show}
	<dialog id="confirm_modal" class="modal" open>
		<div class="modal-box">
			<form method="dialog">
				<button class="btn absolute top-2 right-2 btn-circle btn-ghost btn-sm" onclick={onCancel}>
					✕
				</button>
			</form>

			<div class="mb-4 flex items-center gap-3">
				<div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-error/10">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6 text-error"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
						/>
					</svg>
				</div>
				<div>
					<h3 class="text-lg font-bold">{title}</h3>
					<p class="text-sm text-base-content/70">{message}</p>
				</div>
			</div>

			<div class="modal-action">
				<form method="dialog" class="flex gap-2">
					<button class="btn" onclick={onCancel} disabled={isProcessing}>
						{cancelText}
					</button>
					<button class="btn btn-error" onclick={onConfirm} disabled={isProcessing} type="button">
						{#if isProcessing}
							<span class="loading loading-spinner"></span>
						{:else}
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
								/>
							</svg>
						{/if}
						{confirmText}
					</button>
				</form>
			</div>
		</div>
		<button
			type="button"
			class="modal-backdrop cursor-default"
			onclick={onCancel}
			aria-label="Tutup modal"
		>
		</button>
	</dialog>
{/if}
