<template>
  <div class="card">
    <Loading message="Loading Receipt File" v-if="loading" />
    <div v-else>
      <div class="card-header d-flex justify-content-between">
        <span> Viewing Receipt for {{ receipt.label }}</span>
        <button class="btn btn-sm btn-danger" @click="handleDeleteFile">
          <font-awesome-icon icon="fa-solid fa-trash" />
        </button>
      </div>
      <div class="card-body" v-if="receiptFileData">
        <iframe :src="receiptFileData" width="100%" height="600"></iframe>
      </div>
      <div class="card-body" v-else>
        <div class="mb-3">
          <label for="file">Upload A Receipt</label>
          <input type="file" name="file" id="file" class="form-control" @change="handleFileInput" />
        </div>
        <div class="row">
          <div class="col-lg-4 offset-4">
            <div class="mb-3 text-center d-grid gap-2">
              <button class="btn btn-primary" @click="handleFileUpload" :disabled="!receiptFile">
                UploadFile
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { ref, type Ref, type PropType } from 'vue';

import Loading from '@/components/Loading.vue';
import { ReceiptRequest } from '@/api';
import type { IReceipt } from '@/api/types/receipt';

export default defineComponent({
  components: {
    Loading,
  },
  props: {
    receipt: {
      type: Object as PropType<IReceipt>,
      required: true,
    },
  },
  methods: {
    async handleFileInput(e: Event) {
      const input = e.target as HTMLInputElement;
      if (!input.files?.length) {
        return;
      }

      this.receiptFile = input.files[0];
    },
    async handleFileUpload() {
      if (!this.receiptFile) return;
      this.loading = true;

      ReceiptRequest.PostFile(this.receipt.id, this.receiptFile.type, this.receiptFile).then(
        (r) => {
          ReceiptRequest.GetFile(this.receipt.id).then((r) => {
            this.receiptFileData = URL.createObjectURL(r);
            this.loading = false;
          });
        }
      );
    },
    async handleDeleteFile() {
      this.loading = true;

      ReceiptRequest.DeleteFile(this.receipt.id).then((r) => {
        this.receiptFileData = undefined;
        this.loading = false;
      });
    },
  },

  setup(props) {
    const loading: Ref<Boolean> = ref(true);
    const receiptFileData: Ref<string | undefined> = ref();
    const receiptFile: Ref<File | undefined> = ref();

    Promise.all([
      ReceiptRequest.GetFile(props.receipt.id)
        .then((r) => {
          receiptFileData.value = URL.createObjectURL(r);
          loading.value = false;
        })
        .catch((e) => {
          console.log(e);
          loading.value = false;
        }),
    ]);

    //  try {
    //     ReceiptRequest.GetFile(receiptID)
    //         .then((r) => {
    //           receiptFileData.value = URL.createObjectURL(r);
    //           console.log(receiptFileData.value);
    //         })
    //         .catch((e: Error) => {
    //           console.log(e.message);
    //         }),
    //  } catch(e) {}

    return {
      loading,
      receiptFileData,
      receiptFile,
    };
  },
});
</script>
