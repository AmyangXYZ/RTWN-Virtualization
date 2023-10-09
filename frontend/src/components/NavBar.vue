<template>
  <vs-navbar color="#182A60" text-color="black" class="bar">
    <template #title>
      <vs-navbar-title class="title">
        <a href="/">RTWN Emulator</a>
      </vs-navbar-title>
    </template>
    <vs-navbar-item class="item" index="0">
      <a href="/">Dashboard</a>
    </vs-navbar-item>
    <vs-navbar-item class="item" index="1">
      <a href="/">Applications</a>
    </vs-navbar-item>
    <vs-navbar-item class="item" index="2">
      <a href="/about">About</a>
    </vs-navbar-item>
    <!-- <vs-navbar-item class="item" index="3">
      <a href="/plots">Plots</a>
    </vs-navbar-item> -->
  </vs-navbar>
</template>

<style scoped>
.bar {
  position: absolute;
  width: 100%;
  margin-top: -80px;
  padding-left: 4.16%;
  padding-right: 4.16%;
  height: 105px;
  z-index: -1;
}
.title {
  width: 100%;
  margin-bottom: 25px;
}
.title a {
  color: white;
  font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial,
    sans-serif;
  font-size: 1.6rem;
  font-weight: 300;
}
.item {
  margin-bottom: 25px;
}
.item a {
  color: white;
  font-size: 1.1rem;
  font-weight: 500;
}
</style>
