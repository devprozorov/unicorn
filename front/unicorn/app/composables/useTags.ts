import { ref, computed } from 'vue'
import {
  tagTaxonomy,
  tagsGroups,
  getPrimaryCategories,
  getSpecializationsForCategory,
  getDetailsForSpecialization,
  getChildTags,
  getTagPath,
  getPrimaryCategory,
  searchTags,
  matchesTagFilter,
  getTagsByPriority,
  formatTagDisplay,
  getTagFullDescription,
  type TagOption,
  type TagGroup,
  type TagHierarchy
} from '~/utils/tags'

export function useTags() {
  const selectedCategory = ref<string | null>(null)
  const selectedSpecialization = ref<string | null>(null)
  const selectedTags = ref<string[]>([])
  const searchQuery = ref('')

  // Основные категории (IT, GameDev, Startup)
  const primaryCategories = computed(() => getPrimaryCategories())

  // Специализации для выбранной категории
  const availableSpecializations = computed(() => {
    if (!selectedCategory.value) return []
    return getSpecializationsForCategory(selectedCategory.value)
  })

  // Детальные теги для выбранной специализации
  const availableDetails = computed(() => {
    if (!selectedSpecialization.value) return []
    return getDetailsForSpecialization(selectedSpecialization.value)
  })

  // Все группы тегов с приоритетом
  const tagGroupsByPriority = computed(() => getTagsByPriority())

  // Результаты поиска тегов
  const searchResults = computed(() => {
    if (!searchQuery.value || searchQuery.value.length < 2) return []
    return searchTags(searchQuery.value)
  })

  /**
   * Выбрать категорию
   */
  function selectCategory(category: string) {
    selectedCategory.value = category
    selectedSpecialization.value = null
    
    // Автоматически добавляем категорию в выбранные теги
    if (!selectedTags.value.includes(category)) {
      selectedTags.value.push(category)
    }
  }

  /**
   * Выбрать специализацию
   */
  function selectSpecialization(specialization: string) {
    selectedSpecialization.value = specialization
    
    // Автоматически добавляем специализацию в выбранные теги
    if (!selectedTags.value.includes(specialization)) {
      selectedTags.value.push(specialization)
    }
  }

  /**
   * Добавить/удалить тег
   */
  function toggleTag(tag: string) {
    const index = selectedTags.value.indexOf(tag)
    
    if (index > -1) {
      selectedTags.value.splice(index, 1)
      
      // Если удаляем категорию, сбрасываем связанные выборы
      if (tag === selectedCategory.value) {
        selectedCategory.value = null
        selectedSpecialization.value = null
      }
      
      // Если удаляем специализацию, сбрасываем её выбор
      if (tag === selectedSpecialization.value) {
        selectedSpecialization.value = null
      }
    } else {
      selectedTags.value.push(tag)
      
      // Автоматически выбираем категорию и специализацию на основе пути тега
      const path = getTagPath(tag)
      if (path.length > 0 && path[0] && !selectedCategory.value) {
        selectedCategory.value = path[0]
      }
      if (path.length > 1 && path[1] && !selectedSpecialization.value) {
        selectedSpecialization.value = path[1]
      }
    }
  }

  /**
   * Очистить все выбранные теги
   */
  function clearAllTags() {
    selectedTags.value = []
    selectedCategory.value = null
    selectedSpecialization.value = null
  }

  /**
   * Очистить категорию и связанные теги
   */
  function clearCategory() {
    if (!selectedCategory.value) return
    
    const childTags = getChildTags(selectedCategory.value)
    selectedTags.value = selectedTags.value.filter(
      tag => tag !== selectedCategory.value && !childTags.includes(tag)
    )
    
    selectedCategory.value = null
    selectedSpecialization.value = null
  }

  /**
   * Получить отображаемое имя тега
   */
  function getTagDisplayName(tag: string): string {
    return formatTagDisplay(tag)
  }

  /**
   * Получить полное описание тега для tooltip
   */
  function getTagDescription(tag: string): string {
    return getTagFullDescription(tag)
  }

  /**
   * Проверить, соответствуют ли теги вакансии фильтру
   */
  function checkTagMatch(vacancyTags: string[]): boolean {
    return matchesTagFilter(vacancyTags, selectedTags.value)
  }

  /**
   * Получить количество вакансий, соответствующих фильтру
   */
  function getMatchingCount(vacancies: Array<{ tags?: string[] }>): number {
    return vacancies.filter(v => checkTagMatch(v.tags || [])).length
  }

  /**
   * Быстрая установка категории (для кнопок IT, GameDev, Startup)
   */
  function setQuickCategory(category: 'it' | 'gamedev' | 'startup' | '') {
    clearAllTags()
    
    const categoryMap: Record<string, string> = {
      it: 'IT',
      gamedev: 'GameDev',
      startup: 'Startup'
    }
    
    if (category && categoryMap[category]) {
      selectCategory(categoryMap[category])
    }
  }

  /**
   * Получить все теги для категории (включая дочерние)
   */
  function getAllTagsForCategory(category: string): string[] {
    const tags = [category]
    const children = getChildTags(category)
    return [...tags, ...children]
  }

  /**
   * Экспорт выбранных тегов для API
   */
  function getSelectedTagsForApi(): string[] {
    return selectedTags.value
  }

  /**
   * Импорт тегов из API
   */
  function setTagsFromApi(tags: string[]) {
    selectedTags.value = [...tags]
    
    // Автоматически определяем категорию и специализацию
    if (tags.length > 0 && tags[0]) {
      const firstTagPath = getTagPath(tags[0])
      if (firstTagPath.length > 0 && firstTagPath[0]) {
        selectedCategory.value = firstTagPath[0]
      }
      if (firstTagPath.length > 1 && firstTagPath[1]) {
        selectedSpecialization.value = firstTagPath[1]
      }
    }
  }

  /**
   * Получить рекомендуемые теги на основе уже выбранных
   */
  function getRecommendedTags(): TagOption[] {
    if (selectedTags.value.length === 0 || !selectedTags.value[0]) return []
    
    const recommended: TagOption[] = []
    const primaryCat = getPrimaryCategory(selectedTags.value[0])
    
    if (!primaryCat) return []
    
    // Рекомендуем специализации из той же категории
    const specializations = getSpecializationsForCategory(primaryCat)
    
    for (const spec of specializations) {
      if (!selectedTags.value.includes(spec.value)) {
        recommended.push(spec)
      }
    }
    
    return recommended.slice(0, 5) // Максимум 5 рекомендаций
  }

  return {
    // State
    selectedCategory,
    selectedSpecialization,
    selectedTags,
    searchQuery,
    
    // Computed
    primaryCategories,
    availableSpecializations,
    availableDetails,
    tagGroupsByPriority,
    searchResults,
    
    // Methods
    selectCategory,
    selectSpecialization,
    toggleTag,
    clearAllTags,
    clearCategory,
    getTagDisplayName,
    getTagDescription,
    checkTagMatch,
    getMatchingCount,
    setQuickCategory,
    getAllTagsForCategory,
    getSelectedTagsForApi,
    setTagsFromApi,
    getRecommendedTags,
    
    // Utilities
    tagTaxonomy,
    tagsGroups
  }
}
