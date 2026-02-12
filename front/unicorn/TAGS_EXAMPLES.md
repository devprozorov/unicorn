# Примеры использования новой системы тегов

## 1. Базовое использование в компоненте

```vue
<template>
  <div>
    <HierarchicalTagSelector
      :selected-tags="selectedTags"
      :selected-category="selectedCategory"
      :selected-specialization="selectedSpecialization"
      :enable-search="true"
      :show-recommendations="true"
      @toggle-tag="handleToggleTag"
      @select-category="handleSelectCategory"
      @select-specialization="handleSelectSpecialization"
      @clear-all="handleClearAll"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTags } from '~/composables/useTags'

const tags = useTags()

const selectedTags = ref<string[]>([])
const selectedCategory = ref<string | null>(null)
const selectedSpecialization = ref<string | null>(null)

function handleToggleTag(tag: string) {
  tags.toggleTag(tag)
  selectedTags.value = tags.getSelectedTagsForApi()
}

function handleSelectCategory(category: string) {
  tags.selectCategory(category)
  selectedCategory.value = category
  selectedTags.value = tags.getSelectedTagsForApi()
}

function handleSelectSpecialization(specialization: string) {
  tags.selectSpecialization(specialization)
  selectedSpecialization.value = specialization
  selectedTags.value = tags.getSelectedTagsForApi()
}

function handleClearAll() {
  tags.clearAllTags()
  selectedTags.value = []
  selectedCategory.value = null
  selectedSpecialization.value = null
}
</script>
```

## 2. Фильтрация вакансий с иерархией

```typescript
import { matchesTagFilter } from '~/utils/tags'

// Пример вакансий
const vacancies = [
  { id: 1, title: 'Backend разработчик', tags: ['IT', 'Backend', 'Backend+Node.js'] },
  { id: 2, title: 'Frontend разработчик', tags: ['IT', 'Frontend', 'Frontend+React'] },
  { id: 3, title: 'Unity разработчик', tags: ['GameDev', 'Unity', 'Unity+C#'] },
]

// Фильтр по категории IT
const filterTags = ['IT']
const filtered1 = vacancies.filter(v => matchesTagFilter(v.tags, filterTags))
// Результат: вакансии 1 и 2 (обе в IT)

// Фильтр по Backend
const filterTags2 = ['Backend']
const filtered2 = vacancies.filter(v => matchesTagFilter(v.tags, filterTags2))
// Результат: вакансия 1 (Backend+Node.js является дочерним для Backend)

// Фильтр по конкретной технологии
const filterTags3 = ['Backend+Node.js']
const filtered3 = vacancies.filter(v => matchesTagFilter(v.tags, filterTags3))
// Результат: только вакансия 1
```

## 3. Поиск тегов

```typescript
import { searchTags } from '~/utils/tags'

// Поиск по C#
const results1 = searchTags('C#')
// => ['Backend+C#', 'Unity+C#', 'UE5+C#']

// Поиск по Node
const results2 = searchTags('Node')
// => ['Backend+Node.js']

// Поиск по Unity
const results3 = searchTags('Unity')
// => ['Unity', 'Unity+C#', 'Unity+Mobile', 'Unity+PC']
```

## 4. Работа с путями тегов

```typescript
import { getTagPath, getParentTag, getChildTags } from '~/utils/tags'

// Получить полный путь
const path = getTagPath('Backend+Node.js')
// => ['IT', 'Backend', 'Backend+Node.js']

// Получить родителя
const parent = getParentTag('Backend+Node.js')
// => 'Backend'

// Получить всех потомков
const children = getChildTags('IT')
// => ['Backend', 'Frontend', 'Mobile', ..., 'Backend+Node.js', 'Frontend+React', ...]
```

## 5. Быстрый выбор категории (кнопки IT, GameDev, Startup)

```vue
<template>
  <div class="quick-filters">
    <button 
      @click="setQuickCategory('it')"
      :class="{ active: quickCategory === 'it' }"
    >
      IT 💻
    </button>
    <button 
      @click="setQuickCategory('gamedev')"
      :class="{ active: quickCategory === 'gamedev' }"
    >
      GameDev 🎮
    </button>
    <button 
      @click="setQuickCategory('startup')"
      :class="{ active: quickCategory === 'startup' }"
    >
      Startup 🚀
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTags } from '~/composables/useTags'

const tags = useTags()
const quickCategory = ref('')

function setQuickCategory(category: 'it' | 'gamedev' | 'startup') {
  quickCategory.value = category
  tags.setQuickCategory(category)
}
</script>
```

## 6. Создание вакансии с тегами

```typescript
// При создании вакансии компанией
const tags = useTags()

// Выбираем категорию
tags.selectCategory('GameDev')

// Выбираем специализацию
tags.selectSpecialization('Unity')

// Добавляем детальный тег
tags.toggleTag('Unity+C#')

// Получаем теги для отправки на сервер
const tagsForApi = tags.getSelectedTagsForApi()
// => ['GameDev', 'Unity', 'Unity+C#']

// Отправляем на сервер
await axios.post('/vacancies', {
  title: 'Unity разработчик',
  description: 'Ищем Unity разработчика на C#',
  tags: tagsForApi,
  // ... другие поля
})
```

## 7. Отображение тегов в карточке вакансии

```vue
<template>
  <div class="vacancy-card">
    <h3>{{ vacancy.title }}</h3>
    <div class="tags">
      <span 
        v-for="tag in vacancy.tags" 
        :key="tag"
        class="tag"
        :title="getTagDescription(tag)"
      >
        {{ getTagDisplayName(tag) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTags } from '~/composables/useTags'

const props = defineProps<{
  vacancy: {
    title: string
    tags: string[]
  }
}>()

const tags = useTags()

function getTagDisplayName(tag: string) {
  return tags.getTagDisplayName(tag)
  // 'Backend+Node.js' => 'Node.js'
  // 'Unity+C#' => 'C#'
}

function getTagDescription(tag: string) {
  return tags.getTagDescription(tag)
  // 'Backend+Node.js' => 'IT > Backend > Node.js'
  // 'Unity+C#' => 'GameDev > Unity > C#'
}
</script>
```

## 8. Рекомендации тегов

```typescript
import { useTags } from '~/composables/useTags'

const tags = useTags()

// Пользователь выбрал 'Backend'
tags.selectCategory('IT')
tags.selectSpecialization('Backend')

// Получаем рекомендации
const recommended = tags.getRecommendedTags()
// => [
//   { value: 'Frontend', labelKey: 'tags.specializations.frontend' },
//   { value: 'Mobile', labelKey: 'tags.specializations.mobile' },
//   { value: 'DevOps', labelKey: 'tags.specializations.devops' },
//   ...
// ]
```

## 9. Комбинация тегов из разных категорий

```typescript
// Вакансия может иметь теги из разных сфер
const vacancy = {
  title: 'Full Stack разработчик для игровой платформы',
  tags: [
    'IT',
    'Backend',
    'Backend+Node.js',
    'Frontend',
    'Frontend+React',
    'GameDev' // Также связана с GameDev индустрией
  ]
}

// При фильтрации по 'IT' - вакансия будет показана
// При фильтрации по 'GameDev' - вакансия также будет показана
// При фильтрации по 'Backend' - вакансия будет показана
```

## 10. Миграция существующих вакансий

```typescript
// Старый формат (плоский список)
const oldVacancy = {
  tags: ['Backend', 'Node.js', 'JavaScript', 'Remote']
}

// Новый формат (иерархический)
const newVacancy = {
  tags: [
    'IT',                    // Категория
    'Backend',               // Специализация
    'Backend+Node.js',       // Детальный тег (включает Node.js)
    'Remote'                 // Старый тег (сохраняется для обратной совместимости)
  ]
}

// Функция миграции
function migrateVacancyTags(oldTags: string[]): string[] {
  const newTags: string[] = []
  
  // Определяем категорию
  if (oldTags.includes('Backend') || oldTags.includes('Frontend') || oldTags.includes('Fullstack')) {
    newTags.push('IT')
  }
  
  // Добавляем специализацию и детальные теги
  if (oldTags.includes('Backend')) {
    newTags.push('Backend')
    
    if (oldTags.includes('Node.js')) {
      newTags.push('Backend+Node.js')
    } else if (oldTags.includes('Python')) {
      newTags.push('Backend+Python')
    }
  }
  
  // Сохраняем остальные теги
  const keepTags = oldTags.filter(tag => 
    !['Backend', 'Frontend', 'Node.js', 'Python'].includes(tag)
  )
  
  return [...newTags, ...keepTags]
}
```

## Преимущества новой системы

1. **Иерархическая структура** - понятная навигация от общего к частному
2. **Переиспользование тегов** - C# может быть в Backend+C# и Unity+C#
3. **Умная фильтрация** - выбор родительского тега показывает все дочерние
4. **Поиск** - можно искать по любому уровню иерархии
5. **UX** - пользователь видит сначала основные категории (IT, GameDev), потом уточняет
6. **Масштабируемость** - легко добавлять новые категории и специализации
7. **Обратная совместимость** - старые теги продолжают работать
