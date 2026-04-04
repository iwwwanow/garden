<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';

  const entries = [
    { rank: 1,  name: 'Ada Lovelace',   fdBalance: 2048, daysGrown: 155, waterings: 3102 },
    { rank: 2,  name: 'Marie Curie',     fdBalance: 1024, daysGrown: 120, waterings: 2548 },
    { rank: 3,  name: 'Nikola Tesla',    fdBalance: 512,  daysGrown: 87,  waterings: 1240 },
    { rank: 4,  name: 'Alan Turing',     fdBalance: 256,  daysGrown: 63,  waterings: 900  },
    { rank: 5,  name: 'Emmy Noether',    fdBalance: 128,  daysGrown: 42,  waterings: 630  },
    { rank: 6,  name: 'Rosalind Franklin', fdBalance: 64, daysGrown: 21, waterings: 310  },
    { rank: 7,  name: 'John Doe',        fdBalance: 32,   daysGrown: 7,   waterings: 75   },
  ];

  const medals = ['🥇', '🥈', '🥉'];
</script>

<div class="page">
  <h2 class="title">Таблица</h2>

  <!-- Top-3 podium -->
  <div class="podium">
    {#each entries.slice(0, 3) as e}
      <GlassCard column padding="12px" gap="6px" center>
        <span class="medal">{medals[e.rank - 1]}</span>
        <div class="podium-avatar">
          <span>{e.name[0]}</span>
        </div>
        <span class="podium-name">{e.name.split(' ')[0]}</span>
        <span class="podium-fd">{e.fdBalance} <span class="fd-u">FD</span></span>
      </GlassCard>
    {/each}
  </div>

  <!-- Rest of the list -->
  <div class="list">
    {#each entries.slice(3) as e}
      <GlassCard padding="10px 14px">
        <div class="row">
          <span class="rank">#{e.rank}</span>
          <div class="list-avatar">
            <span>{e.name[0]}</span>
          </div>
          <div class="info">
            <span class="rname">{e.name}</span>
            <div class="rstats">
              <span>{e.daysGrown} дней</span>
              <span class="sep">·</span>
              <span>{e.waterings} поливов</span>
            </div>
          </div>
          <span class="rfd">{e.fdBalance} <span class="fdu">FD</span></span>
        </div>
      </GlassCard>
    {/each}
  </div>
</div>

<style>
.page {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.title {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 32px;
  padding: 8px 4px 0;
  color: var(--dark);
}

/* Podium */
.podium {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 8px;
}

.medal { font-size: 24px; line-height: 1; }

.podium-avatar {
  width: 48px; height: 48px;
  border-radius: 50%;
  background: linear-gradient(180deg, var(--blue-light), var(--blue-mid));
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 500;
  color: white;
  box-shadow: 0 2px 8px rgba(0,120,215,0.3);
}

.podium-name {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 14px;
  color: var(--dark);
  text-align: center;
}

.podium-fd {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 16px;
  color: var(--blue-mid);
}
.fd-u { font-size: 11px; }

/* List */
.list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.row {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.rank {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 16px;
  color: var(--gray-soft);
  min-width: 28px;
}

.list-avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: linear-gradient(180deg, var(--mint-soft), var(--mint));
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 500;
  color: var(--dark);
  flex-shrink: 0;
}

.info { flex: 1; display: flex; flex-direction: column; gap: 1px; }

.rname {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 16px;
  color: var(--dark);
}

.rstats {
  display: flex;
  gap: 4px;
  font-size: 11px;
  color: var(--gray-mid);
}
.sep { color: var(--gray-light); }

.rfd {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 15px;
  color: var(--blue-mid);
  flex-shrink: 0;
}
.fdu { font-size: 10px; }
</style>
