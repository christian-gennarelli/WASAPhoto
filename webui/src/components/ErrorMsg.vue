<script>
    export default {
        props: ['error']
    }
</script>

<template>
    <div v-if="error.ErrorCode">
        <p> {{ error.ErrorCode + " " +  error.Description }}  </p>
    </div>
</template>