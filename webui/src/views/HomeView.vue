<script>
    import Header from '../components/Header.vue'
    import Post from '../components/Post.vue'
    export default {
        data() {
            return{ 
                user: { // User info
                    ID: '',
                    Username: '',
                    Name: '',
                    Birthdate: '',
                    ProfilePic: '',
                },
                posts: [],
                loading: true,
            }
        },
        methods: {
            async getUserStream(){
                await this.$axios.get(
                    '/users/' + this.user.Username + '/stream',
                    { 
                        headers: {
                            'Authorization': this.user.ID,
                        }
                    }
                ).then((res) => {
                    this.posts = res.data 
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            },
        },
        components: {
            Header: Header,
            Post: Post
        },
        created() {

            // Username and token
            this.user.Username = localStorage.getItem('Username')
            this.user.ID = localStorage.getItem('ID')
            this.user.ProfilePic = localStorage.getItem('ProfilePic')
            this.user.Birthdate = localStorage.getItem('Birthdate')
            this.user.Name = localStorage.getItem('Name')

            // Stream
            this.getUserStream();
            this.loading = false
        },
    }
</script>

<template>
    <div v-if="!loading">
        <Header></Header>
        <p class="no-posts" v-if="!posts"> No posts to see yet... start following!</p>
        <div v-else class="home-container"> 
            <Post
                v-if="posts"
                v-for="(post, key) in posts" 
                :post="post" 
                :user="user"
                @add-like="this.posts[key].Likes ? this.posts[key].Likes.push(this.user) :  this.posts[key].Likes = [user]"
                @remove-like="this.posts[key].Likes = this.posts[key].Likes.filter(u => u.Username !== this.user.Username)"
                @remove-post="this.posts.splice(key, 1)"
            ></Post>
        </div>
    </div>
</template>

<style>

.home-container {
    margin: 15px;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-auto-rows: auto;
}

.no-posts {
    text-align: center;
}

</style>